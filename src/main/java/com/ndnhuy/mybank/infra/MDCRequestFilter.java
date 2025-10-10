package com.ndnhuy.mybank.infra;

import jakarta.servlet.*;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.extern.slf4j.Slf4j;
import org.slf4j.MDC;
import org.springframework.core.annotation.Order;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.util.Enumeration;

@Component
@Slf4j
@Order(1)
public class MDCRequestFilter extends OncePerRequestFilter {

  @Override
  protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain filterChain) throws ServletException, IOException {
    try {
      MDC.put("traceId", java.util.UUID.randomUUID().toString());

      String header = "X-Request-Id";
      String reqId = request.getHeader(header);

      filterChain.doFilter(request, response);
    } finally {
      MDC.clear();
    }

  }
}
