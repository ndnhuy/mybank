package com.ndnhuy.mybank.infra;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.Gauge;
import io.micrometer.core.instrument.Timer;
import lombok.Builder;
import lombok.Getter;

@Builder
@Getter
public class QueueMetrics {
  // Arrival metrics
  private final Counter transfersSubmitted;
  private final Timer submissionTime;

  // Service metrics  
  private final Timer serviceTime;
  private final Counter transfersCompleted;

  // Queue metrics
  private final Gauge queueLength;
  private final Timer waitTime;

  // System metrics
  private final Gauge systemUtilization;
}
