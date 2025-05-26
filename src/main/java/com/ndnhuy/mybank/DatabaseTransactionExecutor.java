package com.ndnhuy.mybank;

import org.springframework.stereotype.Component;

import jakarta.transaction.Transactional;

@Component
public class DatabaseTransactionExecutor {

  @Transactional
  public void executeInTransaction(Runnable callback) {
    callback.run();
  }
}
