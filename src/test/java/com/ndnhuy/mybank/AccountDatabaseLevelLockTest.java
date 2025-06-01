package com.ndnhuy.mybank;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.bean.override.mockito.MockitoBean;

import org.junit.jupiter.api.BeforeEach;
import static org.assertj.core.api.AssertionsForClassTypes.assertThat;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@AutoConfigureMockMvc
class AccountDatabaseLevelLockTest {

  private final static String ACCOUNT_ID_PREFIX = "test-account-";

  private final static String generateAccountId() {
    return ACCOUNT_ID_PREFIX + System.currentTimeMillis();
  }

  @Autowired
  private BankService bankService;

  @Autowired
  private AccountRepository accountRepository;

  // Use MockitoBean to disable AccountLockService so that we can test
  // the database-level locking mechanism without interference from the service
  // layer
  @MockitoBean
  private LocalLockService<String> localLockService;

  @BeforeEach
  void setUp() {
    var testAccounts = accountRepository.findAll().stream()
        .filter(account -> account.getId().startsWith(ACCOUNT_ID_PREFIX))
        .toList();
    accountRepository.deleteAll(testAccounts);
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