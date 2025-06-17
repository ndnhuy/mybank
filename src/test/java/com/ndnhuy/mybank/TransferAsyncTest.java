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

    var sourceAccountCount = 10;
    for (int i = 0; i < sourceAccountCount; i++) {
      sourceAccounts.add(bankService.createAccount(generateAccountId(), 100.0));
    }

    var destAccountCount = 10;
    for (int i = 0; i < destAccountCount; i++) {
      destinationAccounts.add(bankService.createAccount(generateAccountId(), 0.0));
    }

    // Concurrent execution setup
    ScheduledExecutorService scheduler = Executors.newScheduledThreadPool(2);
    AtomicInteger requestCount = new AtomicInteger(0);
    var totalRequests = 100; // Total number of requests to submit
    CountDownLatch allRequestsSubmitted = new CountDownLatch(totalRequests);
    List<java.util.concurrent.FutureTask<Void>> futures = new ArrayList<>();

    try {
      // when - submit transfer requests at controlled arrival rate 
      var RPS = 50;
      scheduler.scheduleAtFixedRate(() -> {
        if (requestCount.get() < totalRequests) {
          int idx = requestCount.getAndIncrement();
          Account fromAccount = sourceAccounts.get(idx % sourceAccounts.size());
          Account toAccount = destinationAccounts.get(idx % destinationAccounts.size());

          var future = bankDeskService.submitTransfer(fromAccount.getId(), toAccount.getId(), 1.0);
          synchronized (futures) {
            futures.add(future);
          }

          allRequestsSubmitted.countDown();
        }
      }, 0, 1000 / RPS, TimeUnit.MILLISECONDS);

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
      assertThat(queueMetrics.getTransfersSubmittedCount()).isEqualTo(totalRequests);
      assertThat(queueMetrics.getTransfersCompleted().count()).isEqualTo(totalRequests);
      assertThat(queueMetrics.getWaitTime().totalTime(TimeUnit.MILLISECONDS)).isGreaterThan(0); // Queue buildup occurred
      assertThat(queueMetrics.getServiceTime().totalTime(TimeUnit.MILLISECONDS)).isGreaterThan(0);

      // Verify account balances - each source account should have lost some money
      for (Account sourceAccount : sourceAccounts) {
        var updatedAccount = bankService.getAccountInfo(sourceAccount.getId());
        assertThat(updatedAccount.getBalance()).isLessThan(100.0);
      }

      // Verify destination accounts - each should have received money (> 0)
      for (Account destAccount : destinationAccounts) {
        var updatedAccount = bankService.getAccountInfo(destAccount.getId());
        assertThat(updatedAccount.getBalance()).isGreaterThan(0.0);
      }

      // Verify total money of all accounts remains constant
      double totalSourceBalance = sourceAccounts.stream()
          .mapToDouble(account -> bankService.getAccountInfo(account.getId()).getBalance())
          .sum();
      double totalDestinationBalance = destinationAccounts.stream()
          .mapToDouble(account -> bankService.getAccountInfo(account.getId()).getBalance())
          .sum();
      assertThat(totalSourceBalance + totalDestinationBalance).isEqualTo(sourceAccountCount * 100); // initial total money

      System.out.println("\n=== Concurrent Transfer Test with Arrival Rate ===");
      report.print();
      System.out.printf("Average wait time per transfer: %.2f ms%n",
          queueMetrics.getWaitTime().totalTime(TimeUnit.MILLISECONDS) / queueMetrics.getTransfersCompleted().count());
      System.out.printf("Average service time per transfer: %.2f ms%n",
          queueMetrics.getServiceTime().totalTime(TimeUnit.MILLISECONDS)
              / queueMetrics.getTransfersCompleted().count());
      System.out.printf("Total money transferred: $%.2f%n", totalDestinationBalance - totalSourceBalance);

    } finally {
      if (!scheduler.isShutdown()) {
        scheduler.shutdownNow();
      }
    }
  }

}
