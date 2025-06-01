package com.ndnhuy.mybank;

import jakarta.validation.constraints.Min;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import lombok.Builder;
import lombok.Getter;

@Builder
@Getter
public class TransferRequest {

  @NotBlank
  private String fromAccountId;

  @NotBlank
  private String toAccountId;

  @NotNull
  @Min(0)
  private Double amount;
}
