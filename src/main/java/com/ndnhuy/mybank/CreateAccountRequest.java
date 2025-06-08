package com.ndnhuy.mybank;

import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotNull;
import lombok.Builder;
import lombok.Getter;

@Builder
@Getter
public class CreateAccountRequest {

  @NotNull
  @Min(0)
  private Double initialBalance;

}
