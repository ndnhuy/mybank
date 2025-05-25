package com.ndnhuy.mybank;

import java.util.List;
import java.util.Map;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.locks.ReentrantLock;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Service
public class BankService {

  private static final Map<String, ReentrantLock> ACCOUNT_LOCKS = new ConcurrentHashMap<>();

  private static final ReentrantLock getLock(String accountId) {
    return ACCOUNT_LOCKS.computeIfAbsent(accountId, k -> new ReentrantLock());
  }

  /**
   * Acquires locks for the specified account IDs in a consistent order to prevent
   * deadlocks.
   * 
   * @param accountIds the account IDs for which to acquire locks
   * @return a Runnable that releases the locks when executed
   */
  private static final Runnable acquireLocks(String... accountIds) {
    if (accountIds == null || accountIds.length == 0) {
      throw new IllegalArgumentException("At least one account ID must be provided");
    }

    log.info("Acquiring locks for accounts: {}", (Object[]) accountIds);

    // Sort the account IDs to ensure consistent locking order
    String[] sortedAccountIds = accountIds.clone();
    java.util.Arrays.sort(sortedAccountIds);

    for (String accountId : sortedAccountIds) {
      getLock(accountId).lock();
      log.info("Acquired lock for account: {}", accountId);
    }

    return () -> {
      for (String accountId : sortedAccountIds) {
        getLock(accountId).unlock();
        log.info("Released lock for account: {}", accountId);
      }
    };
  }

  @Autowired
  private AccountRepository accountRepository;

  public Account createAccount(@NonNull Double initialBalance) {
    return createAccount(UUID.randomUUID().toString(), initialBalance);
  }

  public Account createAccount(@NonNull String accountId, @NonNull Double initialBalance) {
    log.info("Creating account with id: {}, initial balance: {}", accountId, initialBalance);
    if (initialBalance < 0) {
      throw new IllegalArgumentException("Initial balance must be non-negative");
    }

    var acc = accountRepository.save(Account.builder()
        .id(accountId)
        .balance(initialBalance)
        .build());

    log.info("Created account with account id: {}, initial balance: {}", acc.getId(), initialBalance);

    return acc;
  }

  public Account getAccount(String accountNumber) {
    log.info("Retrieving account with id: {}", accountNumber);
    if (accountNumber == null || accountNumber.isEmpty()) {
      throw new IllegalArgumentException("Account number must not be null or empty");
    }

    return accountRepository.findById(accountNumber)
        .orElseThrow(() -> new IllegalArgumentException("Account not found with id: " + accountNumber));
  }

  public void transfer(@NonNull String fromAccId, @NonNull String toAccId, Double amount) {
    var unlock = acquireLocks(fromAccId, toAccId);
    try {

      log.info("Transferring {} from account {} to account {}", amount, fromAccId, toAccId);

      if (amount <= 0) {
        throw new IllegalArgumentException("Transfer amount must be positive");
      }
      var fromAccount = getAccount(fromAccId);
      var toAccount = getAccount(toAccId);
      fromAccount.withdraw(amount);
      toAccount.deposit(amount);
      accountRepository.saveAll(List.of(fromAccount, toAccount));

      log.info("Transfer successful: {} from account {} to account {}", amount, fromAccount, toAccount);
    } finally {
      unlock.run();
      log.info("Released locks for accounts {} and {}", fromAccId, toAccId);
    }
  }
}
