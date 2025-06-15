package com.ndnhuy.mybank;

import java.util.List;

import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import com.ndnhuy.mybank.domain.BankService;

import lombok.AllArgsConstructor;

@RestController
@AllArgsConstructor
public class AccountController {

  private final BankService bankService;

  @PostMapping("/accounts")
  public AccountInfo createAccount(@Validated @RequestBody CreateAccountRequest request) {
    return bankService.createAccount(request.getInitialBalance());
  }

  @GetMapping("/accounts/{accountId}")
  public AccountInfo getAccount(@PathVariable String accountId) {
    return bankService.getAccountInfo(accountId);
  }

  @GetMapping("/accounts")
  public List<AccountInfo> getAllAccounts() {
    return bankService.getAllAccounts().stream()
        .map(acc -> AccountInfo.builder()
            .id(acc.getId())
            .balance(acc.getBalance())
            .build())
        .toList();
  }

  @PostMapping("/accounts/transfer")
  public void transfer(@Validated @RequestBody TransferRequest request) {
    bankService.transfer(request.getFromAccountId(), request.getToAccountId(), request.getAmount());
  }

}
