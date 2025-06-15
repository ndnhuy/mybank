package com.ndnhuy.mybank.domain;

import java.util.concurrent.BlockingQueue;
import java.util.concurrent.Callable;
import java.util.concurrent.FutureTask;
import java.util.concurrent.LinkedBlockingDeque;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.ndnhuy.mybank.TransferRequest;
import com.ndnhuy.mybank.infra.QueueMetrics;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.Gauge;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import jakarta.annotation.PostConstruct;

@Service
public class AsyncBankDeskService implements BankDeskService {

  @Autowired
  private BankService bankService;

  @Autowired
  private MeterRegistry registry;

  private BlockingQueue<TransferTask> transferQueue = new LinkedBlockingDeque<>(100);

  // Track if worker is currently processing a task
  private volatile boolean workerBusy = false;

  private QueueMetrics metrics;

  // TransferTask wraps TransferRequest and a FutureTask
  private class TransferTask extends FutureTask<Void> {

    TransferTask(TransferRequest request) {
      this(request, System.nanoTime());
    }

    TransferTask(TransferRequest request, long queuedTime) {
      super(new Callable<Void>() {
        @Override
        public Void call() {
          // Calculate wait time from when task was queued until now (when processing starts)
          long waitTime = System.nanoTime() - queuedTime;
          metrics.getWaitTime().record(waitTime, java.util.concurrent.TimeUnit.NANOSECONDS);

          long startService = System.nanoTime();
          bankService.transfer(request.getFromAccountId(), request.getToAccountId(), request.getAmount());
          long serviceTime = System.nanoTime() - startService;

          metrics.getServiceTime().record(serviceTime, java.util.concurrent.TimeUnit.NANOSECONDS);
          metrics.getTransfersCompleted().increment();
          return null;
        }
      });
    }
  }

  @Override
  public FutureTask<Void> submitTransfer(String fromAccountId, String toAccountId, Double amount) {
    // Record submission processing time (how long it takes to submit)
    Timer.Sample submissionSample = Timer.start(registry);
    metrics.getTransfersSubmitted().increment();

    TransferRequest request = TransferRequest.builder()
        .fromAccountId(fromAccountId)
        .toAccountId(toAccountId)
        .amount(amount)
        .build();
    TransferTask task = new TransferTask(request);
    transferQueue.add(task);

    // Record submission processing time (renamed from arrivalRate for clarity)
    submissionSample.stop(metrics.getSubmissionTime());

    return task;
  }

  @PostConstruct
  public void init() {
    // Initialize metrics after registry is injected
    this.metrics = QueueMetrics.builder()
        .transfersSubmitted(Counter.builder("transfers.submitted").register(registry))
        .submissionTime(Timer.builder("transfers.submission.time").register(registry))
        .serviceTime(Timer.builder("transfers.service.time").register(registry))
        .transfersCompleted(Counter.builder("transfers.completed").register(registry))
        .queueLength(
            Gauge.builder("transfers.queue.length", transferQueue, q -> q.size()).register(registry))
        .waitTime(Timer.builder("transfers.wait.time").register(registry))
        .systemUtilization(
            Gauge.builder("system.utilization", this, service -> service.workerBusy ? 1.0 : 0.0).register(registry))
        .build();

    // Start worker thread
    Thread worker = new Thread(() -> {
      while (true) {
        try {
          TransferTask task = transferQueue.take();
          workerBusy = true;
          try {
            task.run();
          } finally {
            workerBusy = false;
          }
        } catch (InterruptedException e) {
          Thread.currentThread().interrupt();
          break; // Exit the loop if interrupted
        }
      }
    });
    worker.start();
  }

}
