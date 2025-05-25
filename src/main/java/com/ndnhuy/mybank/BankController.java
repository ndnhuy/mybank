package com.ndnhuy.mybank;

import lombok.NonNull;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class BankController {

  @Autowired
  private BankService bankService;

  @PostMapping("/transfer")
  public String transfer(@NonNull String fromAccountNumber, @NonNull String toAccountNumber, @NonNull Double amount) {
    bankService.transfer(fromAccountNumber, toAccountNumber, amount);
    return "Transfer successful";
  }
}
