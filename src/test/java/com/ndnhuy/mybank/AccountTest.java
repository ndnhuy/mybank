package com.ndnhuy.mybank;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;

import com.ndnhuy.mybank.domain.Account;
import com.ndnhuy.mybank.domain.BankService;

import static org.junit.jupiter.api.Assertions.*;

import org.junit.jupiter.api.BeforeEach;
import static org.assertj.core.api.AssertionsForClassTypes.assertThat;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@AutoConfigureMockMvc
class AccountTest {

  private final static String ACCOUNT_ID_PREFIX = "test-account-";

  private final static String generateAccountId() {
    return ACCOUNT_ID_PREFIX + System.currentTimeMillis();
  }

  @Autowired
  private BankService bankService;

  @Autowired
  private AccountRepository accountRepository;

  @BeforeEach
  void setUp() {
    var testAccounts = accountRepository.findAll().stream()
        .filter(account -> account.getId().startsWith(ACCOUNT_ID_PREFIX))
        .toList();
    accountRepository.deleteAll(testAccounts);
  }

  @Test
  void testCreateAccount() {
    // given
    Double initialBalance = 100.0;

    // when
    Account account = bankService.createAccount(generateAccountId(), initialBalance);

    // then
    var actualAccount = bankService.getAccountInfo(account.getId());
    assertNotNull(actualAccount);
    assertEquals(initialBalance, actualAccount.getBalance());
  }

  @Test
  void testTransfer() {
    // given
    Account fromAccount = bankService.createAccount(generateAccountId(), 100.0);
    Account toAccount = bankService.createAccount(generateAccountId(), 0.0);

    // when
    bankService.transfer(fromAccount.getId(), toAccount.getId(), 30.0);

    // then
    var fromAccountAfterTransfer = bankService.getAccountInfo(fromAccount.getId());
    var toAccountAfterTransfer = bankService.getAccountInfo(toAccount.getId());
    assertThat(fromAccountAfterTransfer.getBalance()).isEqualTo(70.0);
    assertThat(toAccountAfterTransfer.getBalance()).isEqualTo(30.0);
  }

  @Test
  void testTransferConcurrently() throws InterruptedException {
    // given
    Account fromAccount = bankService.createAccount(generateAccountId(), 100.0);
    Account toAccount = bankService.createAccount(generateAccountId(), 0.0);

    // when
    Thread thread1 = new Thread(() -> bankService.transfer(fromAccount.getId(), toAccount.getId(), 30.0));
    Thread thread2 = new Thread(() -> bankService.transfer(fromAccount.getId(), toAccount.getId(), 20.0));

    thread1.start();
    thread2.start();

    thread1.join();
    thread2.join();

    // then
    var fromAccountAfterTransfer = bankService.getAccountInfo(fromAccount.getId());
    var toAccountAfterTransfer = bankService.getAccountInfo(toAccount.getId());
    assertThat(fromAccountAfterTransfer.getBalance()).isEqualTo(50.0);
    assertThat(toAccountAfterTransfer.getBalance()).isEqualTo(50.0);
  }

  @Test
  void testTransferConcurrently_shouldNotBeDeadlock() throws InterruptedException {
    // given
    Account fromAccount = bankService.createAccount(generateAccountId(), 100.0);
    Account toAccount = bankService.createAccount(generateAccountId(), 100.0);

    // when
    Thread thread1 = new Thread(() -> bankService.transfer(fromAccount.getId(), toAccount.getId(), 50.0));
    Thread thread2 = new Thread(() -> bankService.transfer(toAccount.getId(), fromAccount.getId(), 50.0));

    thread1.start();
    thread2.start();

    thread1.join();
    thread2.join();

    // then
    var fromAccountAfterTransfer = bankService.getAccountInfo(fromAccount.getId());
    var toAccountAfterTransfer = bankService.getAccountInfo(toAccount.getId());
    assertThat(fromAccountAfterTransfer.getBalance()).isEqualTo(100.0);
    assertThat(toAccountAfterTransfer.getBalance()).isEqualTo(100.0);
  }
}