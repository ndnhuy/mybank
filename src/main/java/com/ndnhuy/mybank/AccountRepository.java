package com.ndnhuy.mybank;

import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Lock;
import org.springframework.data.jpa.repository.Query;

import com.ndnhuy.mybank.domain.Account;

import jakarta.persistence.LockModeType;

public interface AccountRepository extends JpaRepository<Account, String> {

  @Lock(LockModeType.PESSIMISTIC_WRITE)
  @Query("SELECT a FROM Account a WHERE a.id = :id")
  Optional<Account> findByIdForUpdate(String id);

}