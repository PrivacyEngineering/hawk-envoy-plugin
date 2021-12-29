use log::{info, trace};
use primes::is_prime;
use proxy_wasm::traits::{Context, HttpContext};
use proxy_wasm::types::{Action, LogLevel};

#[no_mangle]
pub fn _start() {
    proxy_wasm::set_log_level(LogLevel::Trace);
    proxy_wasm::set_http_context(|context_id, _| -> Box<dyn HttpContext> {
        Box::new(PrimeAuthorizer { context_id })
    });
}

struct PrimeAuthorizer {
    context_id: u32,
}

impl Context for PrimeAuthorizer {}

impl HttpContext for PrimeAuthorizer {
    fn on_http_request_headers(&mut self, _: usize) -> Action {
        info!("PrimeAuthorizer started. #{}", self.context_id);

        match self.get_http_request_header("x-prime-token") {
            Some(token) if token.parse::<u64>().is_ok() && is_prime(token.parse().unwrap()) => {
                info!("PrimeAuthorizer. OK. #{} -> Header x-prime-token: {}", self.context_id, token);
                self.resume_http_request();
                Action::Continue
            }
            _ => {
                info!("PrimeAuthorizer. #{} -> Access forbidden", self.context_id);
                self.send_http_response(
                    403,
                    vec![("Powered-By", "proxy-wasm")],
                    Some(b"Access forbidden.\n"),
                );
                Action::Pause
            }
        }
    }
}

