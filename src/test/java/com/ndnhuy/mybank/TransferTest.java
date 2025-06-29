package com.ndnhuy.mybank;

import static org.assertj.core.api.AssertionsForClassTypes.assertThat;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;

import com.ndnhuy.mybank.domain.Account;
import com.ndnhuy.mybank.domain.BankService;
import com.ndnhuy.mybank.domain.DefaultBankDeskService;
import com.ndnhuy.mybank.infra.QueueMetrics;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@AutoConfigureMockMvc
class TransferTest {

  private final static String ACCOUNT_ID_PREFIX = "test-account-";

  private final static String generateAccountId() {
    return ACCOUNT_ID_PREFIX + System.currentTimeMillis();
  }

  @Autowired
  private BankService bankService;

  @Autowired
  private DefaultBankDeskService bankDeskService;

  @Autowired
  private AccountRepository accountRepository;

  @Autowired
  private QueueMetrics queueMetrics;

  @BeforeEach
  void setUp() {
    var testAccounts = accountRepository.findAll().stream()
        .filter(account -> account.getId().startsWith(ACCOUNT_ID_PREFIX))
        .toList();
    accountRepository.deleteAll(testAccounts);
    queueMetrics.reset();
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
    bankDeskService.submitTransfer(fromAccount.getId(), toAccount.getId(), 30.0);

    // then
    var fromAccountAfterTransfer = bankService.getAccountInfo(fromAccount.getId());
    var toAccountAfterTransfer = bankService.getAccountInfo(toAccount.getId());
    assertThat(fromAccountAfterTransfer.getBalance()).isEqualTo(70.0);
    assertThat(toAccountAfterTransfer.getBalance()).isEqualTo(30.0);

    // print metrics
    queueMetrics.getReport().print();
  }

  @Test
  void testTransferConcurrently() throws InterruptedException {
    // given
    Account fromAccount = bankService.createAccount(generateAccountId(), 100.0);
    Account toAccount = bankService.createAccount(generateAccountId(), 0.0);

    // when
    Thread thread1 = new Thread(() -> bankDeskService.submitTransfer(fromAccount.getId(), toAccount.getId(), 30.0));
    Thread thread2 = new Thread(() -> bankDeskService.submitTransfer(fromAccount.getId(), toAccount.getId(), 20.0));

    thread1.start();
    thread2.start();

    thread1.join();
    thread2.join();

    // then
    var fromAccountAfterTransfer = bankService.getAccountInfo(fromAccount.getId());
    var toAccountAfterTransfer = bankService.getAccountInfo(toAccount.getId());
    assertThat(fromAccountAfterTransfer.getBalance()).isEqualTo(50.0);
    assertThat(toAccountAfterTransfer.getBalance()).isEqualTo(50.0);

    queueMetrics.getReport().print();
  }

  @Test
  void testTransferConcurrently_shouldNotBeDeadlock() throws InterruptedException {
    // given
    Account fromAccount = bankService.createAccount(generateAccountId(), 100.0);
    Account toAccount = bankService.createAccount(generateAccountId(), 100.0);

    // when
    Thread thread1 = new Thread(() -> bankDeskService.submitTransfer(fromAccount.getId(), toAccount.getId(), 50.0));
    Thread thread2 = new Thread(() -> bankDeskService.submitTransfer(toAccount.getId(), fromAccount.getId(), 50.0));

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