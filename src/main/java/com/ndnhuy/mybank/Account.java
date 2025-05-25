package com.ndnhuy.mybank;

import jakarta.persistence.*;
import lombok.*;

@Entity
@Table(name = "accounts")
@Getter
@Builder
@AllArgsConstructor
@NoArgsConstructor
@ToString
public class Account {

  @Id
  private String id;

  private Double balance;

  public void deposit(Double amount) {
    if (amount > 0) {
      this.balance += amount;
    } else {
      throw new IllegalArgumentException("Deposit amount must be positive");
    }
  }

  public void withdraw(Double amount) {
    if (amount > 0 && amount <= this.balance) {
      this.balance -= amount;
    } else {
      throw new IllegalArgumentException("Withdrawal amount must be positive and less than or equal to the balance");
    }
  }

}
