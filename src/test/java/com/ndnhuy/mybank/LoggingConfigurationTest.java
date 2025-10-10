package com.ndnhuy.mybank;

import static org.assertj.core.api.Assertions.assertThat;

import java.net.URL;
import java.util.Objects;

import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;
import org.springframework.boot.test.system.CapturedOutput;
import org.springframework.boot.test.system.OutputCaptureExtension;

import ch.qos.logback.classic.LoggerContext;
import ch.qos.logback.classic.joran.JoranConfigurator;
import ch.qos.logback.core.joran.spi.JoranException;

@ExtendWith(OutputCaptureExtension.class)
class LoggingConfigurationTest {

  private static final Logger LOGGER = LoggerFactory.getLogger(LoggingConfigurationTest.class);

  @BeforeAll
  static void loadLogbackConfiguration() throws JoranException {
    LoggerContext context = (LoggerContext) LoggerFactory.getILoggerFactory();
    context.reset();

    JoranConfigurator configurator = new JoranConfigurator();
    configurator.setContext(context);
    URL configurationUrl = Objects.requireNonNull(
        LoggingConfigurationTest.class.getResource("/logback-spring.xml"),
        "logback configuration must be present");
    configurator.doConfigure(configurationUrl);
  }

  @Test
  void shouldOutputStructuredLog_whenMessageIsLogged(CapturedOutput output) {
    // When: a log message is emitted
    String message = "customer created";
    MDC.put("traceId", "test-trace-id");
    LOGGER.info(message);

    // Then: the log contains structured key-value fields
    String logOutput = output.getOut() + output.getErr();
    assertThat(logOutput)
        .contains("[main]")
        .contains("INFO")
        .contains("LoggingConfigurationTest")
        .contains(message);
  }
}
