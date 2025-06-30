package com.ndnhuy.mybank;

import static org.assertj.core.api.Assertions.assertThat;
import static org.springframework.test.web.servlet.result.MockMvcResultHandlers.print;

import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.HttpHeaders;
import org.springframework.test.web.servlet.MockMvc;
import org.springframework.test.web.servlet.request.MockMvcRequestBuilders;
import org.springframework.test.web.servlet.result.MockMvcResultMatchers;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@AutoConfigureMockMvc
public class WebApplicationTest {

  @Autowired
  private MockMvc mockMvc;

  @Test
  void shouldIncludeAllRequestHeadersInResponseHeaders_whenGetAccounts() throws Exception {
    // Given: Custom headers for the request
    String header1 = "X-Test-Header-1";
    String value1 = "value1";
    String header2 = "X-Test-Header-2";
    String value2 = "value2";

    // When: GET /accounts is called with custom headers
    var result = mockMvc.perform(MockMvcRequestBuilders.get("/accounts")
        .header(header1, value1)
        .header(header2, value2))
        .andExpect(MockMvcResultMatchers.status().isOk())
        .andReturn();

    // Then: Response headers should include all request headers
    HttpHeaders responseHeaders = new HttpHeaders();
    result.getResponse().getHeaderNames().forEach(
        h -> responseHeaders.add(h, result.getResponse().getHeader(h)));
    assertThat(responseHeaders.get(header1)).contains(value1);
    assertThat(responseHeaders.get(header2)).contains(value2);
  }

  @Test
  void shouldCreateAccount_whenPostAccountsWithValidInitialBalance() throws Exception {
    // Given: a valid initial balance
    int initialBalance = 1000;
    String requestBody = "{\"initialBalance\":" + initialBalance + "}";

    // When: POST /accounts is called
    var result = mockMvc.perform(MockMvcRequestBuilders.post("/accounts")
        .contentType("application/json")
        .content(requestBody))
        .andDo(print())
        .andExpect(MockMvcResultMatchers.status().isOk())
        .andReturn();

    // Then: Response should contain account info with correct balance
    String response = result.getResponse().getContentAsString();
    assertThat(response).contains("\"balance\":" + initialBalance);
    assertThat(response).contains("\"id\":"); // id should be present
  }
}
