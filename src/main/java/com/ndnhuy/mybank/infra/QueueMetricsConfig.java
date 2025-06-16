package com.ndnhuy.mybank.infra;

import io.micrometer.core.instrument.MeterRegistry;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class QueueMetricsConfig {
    @Bean
    public QueueMetrics queueMetrics(MeterRegistry registry) {
        return new QueueMetrics(registry);
    }
}
