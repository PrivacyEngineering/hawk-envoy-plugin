use log::{info, trace};
use primes::is_prime;
use proxy_wasm::traits::{Context, HttpContext};
use proxy_wasm::types::{Action, Bytes, LogLevel};

#[no_mangle]
pub fn _start() {
    proxy_wasm::set_log_level(LogLevel::Trace);
    proxy_wasm::set_http_context(|context_id, root_context_id| -> Box<dyn HttpContext> {
        Box::new(ActixFilter { context_id, root_context_id })
    });
}

struct ActixFilter {
    context_id: u32,
    root_context_id: u32,
}

impl Context for ActixFilter {
}

impl HttpContext for ActixFilter {

    /// Called when HTTP request headers are received from the client.
    /// Headers can be retrieved using proxy_get_map and/or proxy_get_map_value.
    ///
    /// # Arguments
    /// * i32 (uint32_t) context_id
    /// * i32 (size_t) num_headers
    /// * i32 (bool) end_of_stream
    ///
    /// # Return
    /// * i32 (proxy_action_t) next_action
    fn on_http_request_headers(&mut self, _num_headers: usize) -> Action {
        info!("ActixFilter started. #{}", self.context_id);

        // let headers = self.get_http_request_headers(&self);
        // info!("ActixFilter headers. #{}", headers);

        return Action::Continue
    }

    /// Called for each chunk of HTTP request body received from the client.
    /// Request body can be retrieved using proxy_get_buffer.
    ///
    /// # Arguments
    /// * i32 (uint32_t) context_id
    /// * i32 (size_t) body_size
    /// * i32 (bool) end_of_stream
    ///
    /// # Return
    /// * i32 (proxy_action_t) next_action
    fn on_http_request_body(&mut self, _body_size: usize, _end_of_stream: bool) -> Action {
        // self.get_http_request_body(start, max);
        todo!()
    }

    /// Called when HTTP response headers are received from the upstream.
    /// Headers can be retrieved using proxy_get_map and/or proxy_get_map_value.
    ///
    /// # Arguments
    /// * i32 (uint32_t) context_id
    /// * i32 (size_t) num_headers
    /// * i32 (bool) end_of_stream
    ///
    /// # Return
    /// * i32 (proxy_action_t) next_action
    fn on_http_response_headers(&mut self, _num_headers: usize) -> Action {
        // let headers = self.get_http_response_headers(&self);
        todo!()
    }

    /// Called for each chunk of HTTP response body received from the client.
    /// Response body can be retrieved using proxy_get_buffer.
    ///
    /// # Arguments
    /// * i32 (uint32_t) context_id
    /// * i32 (size_t) body_size
    /// * i32 (bool) end_of_stream
    ///
    /// # Return
    /// * i32 (proxy_action_t) next_action
    fn on_http_response_body(&mut self, _body_size: usize, _end_of_stream: bool) -> Action {
        todo!()
    }

}

