package com.ndnhuy.mybank;

import static org.assertj.core.api.AssertionsForClassTypes.assertThat;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;

import com.ndnhuy.mybank.domain.Account;
import com.ndnhuy.mybank.domain.AsyncBankDeskService;
import com.ndnhuy.mybank.domain.BankService;
import com.ndnhuy.mybank.infra.QueueMetrics;

import lombok.SneakyThrows;

import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.List;
import java.util.ArrayList;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@AutoConfigureMockMvc
public class TransferAsyncTest {

  private final static String ACCOUNT_ID_PREFIX = "test-account-";

  private final static String generateAccountId() {
    return ACCOUNT_ID_PREFIX + System.currentTimeMillis();
  }

  @Autowired
  private BankService bankService;

  @Autowired
  private AsyncBankDeskService bankDeskService;

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
  @SneakyThrows
  void testTransfer() {
    // given
    Account fromAccount = bankService.createAccount(generateAccountId(), 100.0);
    Account toAccount = bankService.createAccount(generateAccountId(), 0.0);

    // when
    var f = bankDeskService.submitTransfer(fromAccount.getId(), toAccount.getId(), 30.0);
    f.get();

    // then
    var fromAccountAfterTransfer = bankService.getAccountInfo(fromAccount.getId());
    var toAccountAfterTransfer = bankService.getAccountInfo(toAccount.getId());
    assertThat(fromAccountAfterTransfer.getBalance()).isEqualTo(70.0);
    assertThat(toAccountAfterTransfer.getBalance()).isEqualTo(30.0);

    queueMetrics.getReport().print();
  }

  @Test
  @SneakyThrows
  void testConcurrentTransfersWithArrivalRate_withDeterministicPeriodicArrivals() {
    // given - setup multiple accounts
    List<Account> sourceAccounts = new ArrayList<>();
    List<Account> destinationAccounts = new ArrayList<>();

    // Create 10 source accounts with $100 each
    for (int i = 0; i < 10; i++) {
      sourceAccounts.add(bankService.createAccount(generateAccountId(), 100.0));
    }

    // Create 10 destination accounts with $0 each
    for (int i = 0; i < 10; i++) {
      destinationAccounts.add(bankService.createAccount(generateAccountId(), 0.0));
    }

    // Concurrent execution setup
    ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(2);
    AtomicInteger requestCount = new AtomicInteger(0);
    CountDownLatch allRequestsSubmitted = new CountDownLatch(25); // 25 total requests
    List<java.util.concurrent.FutureTask<Void>> futures = new ArrayList<>();

    try {
      // when - submit transfer requests at controlled arrival rate (every 100ms = 10 requests/second)
      scheduler.scheduleAtFixedRate(() -> {
        if (requestCount.get() < 25) {
          int idx = requestCount.getAndIncrement();
          Account fromAccount = sourceAccounts.get(idx % sourceAccounts.size());
          Account toAccount = destinationAccounts.get(idx % destinationAccounts.size());

          var future = bankDeskService.submitTransfer(fromAccount.getId(), toAccount.getId(), 1.0);
          synchronized (futures) {
            futures.add(future);
          }

          allRequestsSubmitted.countDown();
        }
      }, 0, 100, TimeUnit.MILLISECONDS);

      // Wait for all requests to be submitted (with timeout)
      boolean submitted = allRequestsSubmitted.await(5, TimeUnit.SECONDS);
      assertThat(submitted).isTrue();

      // Stop scheduling new requests
      scheduler.shutdown();

      // Wait for all transfers to complete
      synchronized (futures) {
        for (var future : futures) {
          future.get(10, TimeUnit.SECONDS);
        }
      }

      // then - verify system behavior and collect metrics
      var report = queueMetrics.getReport();

      // Assert queuing system characteristics
      assertThat(queueMetrics.getTransfersSubmitted().count()).isEqualTo(25);
      assertThat(queueMetrics.getTransfersCompleted().count()).isEqualTo(25);
      assertThat(queueMetrics.getWaitTime().totalTime(TimeUnit.MILLISECONDS)).isGreaterThan(0); // Queue buildup occurred
      assertThat(queueMetrics.getServiceTime().totalTime(TimeUnit.MILLISECONDS)).isGreaterThan(0);

      // Verify account balances - each source account should have lost some money
      for (Account sourceAccount : sourceAccounts) {
        var updatedAccount = bankService.getAccountInfo(sourceAccount.getId());
        assertThat(updatedAccount.getBalance()).isLessThan(100.0);
      }

      // Verify destination accounts received money
      double totalReceived = 0;
      for (Account destAccount : destinationAccounts) {
        var updatedAccount = bankService.getAccountInfo(destAccount.getId());
        totalReceived += updatedAccount.getBalance();
      }
      assertThat(totalReceived).isEqualTo(25.0); // 25 transfers * $1 each

      System.out.println("\n=== Concurrent Transfer Test with Arrival Rate ===");
      report.print();
      System.out.printf("Average wait time per transfer: %.2f ms%n",
          queueMetrics.getWaitTime().totalTime(TimeUnit.MILLISECONDS) / queueMetrics.getTransfersCompleted().count());
      System.out.printf("Average service time per transfer: %.2f ms%n",
          queueMetrics.getServiceTime().totalTime(TimeUnit.MILLISECONDS)
              / queueMetrics.getTransfersCompleted().count());
      System.out.printf("Total money transferred: $%.2f%n", totalReceived);

    } finally {
      if (!scheduler.isShutdown()) {
        scheduler.shutdownNow();
      }
    }
  }

}
