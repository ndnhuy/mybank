package com.ndnhuy.mybank;

import java.io.IOException;
import java.util.Enumeration;

import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

@Component
public class CopyRequestHeadersFilter extends OncePerRequestFilter {
  @Override
  protected void doFilterInternal(
      HttpServletRequest request,
      HttpServletResponse response,
      FilterChain filterChain) throws ServletException, IOException {
    // Only copy 'X-Request-Id' header from request to response
    String header = "X-Request-Id";
    Enumeration<String> values = request.getHeaders(header);
    while (values.hasMoreElements()) {
      response.addHeader(header, values.nextElement());
    }
    filterChain.doFilter(request, response);
  }
}