package com.ndnhuy.mybank.domain;

import java.util.concurrent.Future;

public interface BankDeskService {

  /**
   * Submits a transfer request to the transaction queue.
   * This method is designed to handle the transfer of funds between two accounts
   *
   * @param fromAccountId
   * @param toAccountId
   * @param amount
   */
  Future<Void> submitTransfer(String fromAccountId, String toAccountId, Double amount);
}
