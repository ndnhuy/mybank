package com.ndnhuy.mybank;

import lombok.Builder;
import lombok.Getter;

@Builder
@Getter
public class AccountInfo {
  private String id;
  private Double balance;
}
