package com.ndnhuy.mybank.domain;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Future;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.ndnhuy.mybank.infra.QueueMetrics;

import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class DefaultBankDeskService implements BankDeskService {

  private final BankService bankService;

  @Autowired
  private QueueMetrics metrics;

  @Override
  public Future<Void> submitTransfer(String fromAccountId, String toAccountId, Double amount) {
    metrics.getTransfersSubmitted().increment();
    long start = System.nanoTime();
    bankService.transfer(fromAccountId, toAccountId, amount);
    long duration = System.nanoTime() - start;
    metrics.getServiceTime().record(duration, java.util.concurrent.TimeUnit.NANOSECONDS);
    metrics.getTransfersCompleted().increment();
    return CompletableFuture.completedFuture(null);
  }

}
