package com.ndnhuy.mybank.infra;

import java.util.HashMap;
import java.util.Map;
import java.util.function.Function;

import lombok.extern.slf4j.Slf4j;

/**
 * Utility class to fetch data in a consistent order based on provided keys.
 * It ensures that the keys are sorted and duplicates are handled gracefully.
 * 
 * Initial usecase is to fetch data for accounts in a consistent order to prevent deadlocks.
 */
@Slf4j
public class OrderedKeyDataFetcher {

  public static <T> Map<String, T> fetchDataInOrderedKey(Function<String, T> dataFetcher, String... keys) {
    if (keys == null || keys.length == 0) {
      throw new IllegalArgumentException("At least one key must be provided");
    }

    // Sort keys to ensure consistent ordering
    var sortedKeys = java.util.Arrays.stream(keys)
        .filter(java.util.Objects::nonNull)
        .distinct()
        .sorted()
        .toList();

    if (sortedKeys.size() != keys.length) {
      log.warn("Duplicate keys detected");
    }

    var result = new HashMap<String, T>();
    for (String key : sortedKeys) {
      result.put(key, dataFetcher.apply(key));
    }

    return result;
  }

}
