package com.ndnhuy.mybank.infra;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;
import java.util.Map;
import java.util.Objects;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.locks.ReentrantLock;
import org.springframework.stereotype.Service;

import lombok.extern.slf4j.Slf4j;

/**
 * Generic service for managing local locks on resources identified by keys to prevent concurrent modifications.
 * This service provides a way to acquire locks on multiple resources in a consistent
 * order, ensuring that deadlocks do not occur.
 * 
 * IMPORTANT: This implementation only provides local (in-memory) locking within a 
 * single JVM instance. It does NOT support distributed locking across multiple 
 * application instances or nodes. For distributed systems, consider using external 
 * coordination services like Redis, Zookeeper, or database-based locking mechanisms.
 * 
 * @param <K> the type of keys used to identify resources
 */
@Service
@Slf4j
public class LocalLockService<K extends Comparable<K>> {

  private final Map<K, ReentrantLock> lockMap = new ConcurrentHashMap<>();

  private ReentrantLock getLock(K key) {
    return lockMap.computeIfAbsent(key, k -> new ReentrantLock());
  }

  /**
   * Acquires locks for the specified resource keys in a consistent order to prevent deadlocks.
   * Handles exceptions gracefully by releasing any acquired locks before re-throwing.
   * 
   * @param keys the resource keys for which to acquire locks
   * @return a Runnable that releases the locks when executed
   * @throws IllegalArgumentException if no keys provided or if duplicates exist
   */
  @SafeVarargs
  public final Runnable acquireLocks(K... keys) {
    if (keys == null || keys.length == 0) {
      throw new IllegalArgumentException("At least one key must be provided");
    }

    // Remove duplicates and sort to ensure consistent ordering
    var sortedKeys = Arrays.stream(keys)
        .filter(Objects::nonNull)
        .distinct()
        .sorted() // Sort to ensure consistent order
        .toList();

    if (sortedKeys.size() != keys.length) {
      log.warn("Duplicate keys detected in lock acquisition request");
    }

    log.info("Acquiring locks for resources (sorted): {}", sortedKeys);

    List<K> lockedKeys = new ArrayList<>();
    Runnable cleanUp = () -> {
      Collections.reverse(lockedKeys);
      for (K key : lockedKeys) {
        try {
          getLock(key).unlock();
          log.debug("Released lock for resource: {}", key);
        } catch (Exception e) {
          log.error("Error releasing lock for resource: {}", key, e);
        }
      }
      log.debug("Released all locks for resources: {}", sortedKeys);
    };

    try {
      // Acquire locks one by one, keeping track of what we've acquired
      for (K key : sortedKeys) {
        ReentrantLock lock = getLock(key);
        lock.lock();
        lockedKeys.add(key);
        log.debug("Acquired lock for resource: {}", key);
      }

      log.debug("Successfully acquired all locks for resources: {}", sortedKeys);

      // Return cleanup function that releases locks in reverse order
      return cleanUp;
    } catch (Exception e) {
      // If any lock acquisition fails, release all previously acquired locks
      log.error("Failed to acquire all locks, releasing {} already acquired locks", lockedKeys.size());
      cleanUp.run();
      throw e;
    }
  }
}
