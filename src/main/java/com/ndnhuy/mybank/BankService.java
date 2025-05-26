package com.ndnhuy.mybank;

import java.util.List;
import java.util.UUID;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import jakarta.transaction.Transactional;
import lombok.NonNull;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Service
public class BankService {

  @Autowired
  private LocalLockService<String> localLockService;

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

  public Account getAccountForUpdate(String accountNumber) {
    log.info("Retrieving account with id: {} for update", accountNumber);
    if (accountNumber == null || accountNumber.isEmpty()) {
      throw new IllegalArgumentException("Account number must not be null or empty");
    }

    return accountRepository.findByIdForUpdate(accountNumber)
        .orElseThrow(() -> new IllegalArgumentException("Account not found with id: " + accountNumber));
  }

  @Transactional
  public void transfer(@NonNull String fromAccId, @NonNull String toAccId, Double amount) {
    var unlock = localLockService.acquireLocks(fromAccId, toAccId);
    try {

      log.info("Transferring {} from account {} to account {}", amount, fromAccId, toAccId);

      if (amount <= 0) {
        throw new IllegalArgumentException("Transfer amount must be positive");
      }

      Account fromAccount, toAccount;
      if (fromAccId.compareTo(toAccId) < 0) {
        fromAccount = getAccountForUpdate(fromAccId);
        toAccount = getAccountForUpdate(toAccId);
      } else {
        toAccount = getAccountForUpdate(toAccId);
        fromAccount = getAccountForUpdate(fromAccId);
      }

      fromAccount.withdraw(amount);
      toAccount.deposit(amount);
      accountRepository.saveAll(List.of(fromAccount, toAccount));

      log.info("Transfer successful: {} from account {} to account {}", amount, fromAccount, toAccount);
    } finally {
      if (unlock != null) {
        unlock.run();
        log.info("Releasing locks for accounts: {}, {}", fromAccId, toAccId);
      }
    }
  }
}
