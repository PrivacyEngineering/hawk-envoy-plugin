use std::collections::HashMap;
use actix_web::{App, HttpServer, post, web, Result};
use serde::{Serialize, Deserialize};
use serde_json::Value;

#[derive(Serialize, Deserialize, Debug)]
struct Collect {
    from: String,
    #[serde(flatten)]
    content: HashMap<String, Value>,
}

#[post("/echo")]
async fn echo(info: web::Json<Collect>) -> Result<String> {

    let serialized = serde_json::to_string(&info.content).unwrap();
    println!("received: {}: {}", info.from, serialized);

    Ok(format!("received: {}", info.from))
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new()
            .service(echo)
    })
        .bind("0.0.0.0:8080")?
        .run()
        .await
}