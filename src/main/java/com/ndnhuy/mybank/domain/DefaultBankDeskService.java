package com.ndnhuy.mybank.domain;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Future;

import org.springframework.stereotype.Service;

import lombok.RequiredArgsConstructor;

@Service
@RequiredArgsConstructor
public class DefaultBankDeskService implements BankDeskService {

  private final BankService bankService;

  @Override
  public Future<Void> submitTransfer(String fromAccountId, String toAccountId, Double amount) {
    bankService.transfer(fromAccountId, toAccountId, amount);
    return CompletableFuture.completedFuture(null);
  }

}
