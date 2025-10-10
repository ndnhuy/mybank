package com.ndnhuy.mybank.infra;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.extern.slf4j.Slf4j;
import org.springframework.core.annotation.Order;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.util.Enumeration;

@Component
@Slf4j
@Order(2)
public class RequestLoggingFilter extends OncePerRequestFilter {

  @Override
  protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain filterChain) throws ServletException, IOException {
    String method = request.getMethod();
    String query = request.getQueryString();
    String url = request.getRequestURI() + (query != null ? "?" + query : "");

    StringBuilder headers = new StringBuilder();
    Enumeration<String> headerNames = request.getHeaderNames();
    if (headerNames != null) {
      while (headerNames.hasMoreElements()) {
        String name = headerNames.nextElement();
        String value = request.getHeader(name);
        headers.append(name).append("=").append(value).append("; ");
      }
    }

    log.info("Incoming request: {} {} headers=[{}]", method, url, headers.toString());

    filterChain.doFilter(request, response);
  }
}