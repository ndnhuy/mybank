package com.ndnhuy.mybank;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.LinkedHashSet;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.Set;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.locks.ReentrantLock;
import java.util.stream.Collectors;

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
   * Handles exceptions gracefully by releasing any acquired locks before
   * re-throwing.
   * 
   * @param accountIds the account IDs for which to acquire locks
   * @return a Runnable that releases the locks when executed
   * @throws IllegalArgumentException if no account IDs provided or if duplicates
   *                                  exist
   */
  private static Runnable acquireLocks(String... accountIds) {
    if (accountIds == null || accountIds.length == 0) {
      throw new IllegalArgumentException("At least one account ID must be provided");
    }

    // Remove duplicates and sort to ensure consistent ordering
    Set<String> uniqueAccountIds = Arrays.stream(accountIds)
        .filter(Objects::nonNull)
        .collect(Collectors.toCollection(LinkedHashSet::new));

    if (uniqueAccountIds.size() != accountIds.length) {
      log.warn("Duplicate account IDs detected in lock acquisition request");
    }

    var sortedAccountIds = uniqueAccountIds.stream().sorted().toList();

    log.info("Acquiring locks for accounts (sorted): {}", sortedAccountIds);

    List<String> lockedAccountIds = new ArrayList<>();
    Runnable cleanUp = () -> {
      Collections.reverse(lockedAccountIds);
      for (String accountId : lockedAccountIds) {
        try {
          getLock(accountId).unlock();
          log.debug("Released lock for account: {}", accountId);
        } catch (Exception e) {
          log.error("Error releasing lock for account: {}", accountId, e);
        }
      }
      log.debug("Released all locks for accounts: {}", sortedAccountIds);
    };

    try {
      // Acquire locks one by one, keeping track of what we've acquired
      for (String accountId : sortedAccountIds) {
        ReentrantLock lock = getLock(accountId);
        lock.lock();
        lockedAccountIds.add(accountId);
        log.debug("Acquired lock for account: {}", accountId);
      }

      log.debug("Successfully acquired all locks for accounts: {}", sortedAccountIds);

      // Return cleanup function that releases locks in reverse order
      return cleanUp;
    } catch (Exception e) {
      // If any lock acquisition fails, release all previously acquired locks
      log.error("Failed to acquire all locks, releasing {} already acquired locks", lockedAccountIds.size());
      cleanUp.run();
      throw e;
    }
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
    }
  }
}
