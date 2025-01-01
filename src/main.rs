use tonic::{transport::Channel, Request};

pub mod health {
    tonic::include_proto!("grpc.health.v1");
}

use health::health_client::HealthClient;
use health::HealthCheckRequest;
use health::health_check_response::ServingStatus;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // Connect to the server
    let channel = Channel::from_static("http://localhost:50051")
        .connect()
        .await?;

    // Create the client
    let mut client = HealthClient::new(channel);

    // Create the health check request
    let request = Request::new(HealthCheckRequest {
        service: "".to_string(), // Empty string checks the overall health
    });

    // Send the health check request
    match client.check(request).await {
        Ok(response) => {
            let status = response.into_inner().status;
            if status != ServingStatus::Serving as i32 {
                println!("Service health status: {:?}", status);
                std::process::exit(1);
            }
        }
        Err(e) => {
            println!("Health check failed: {:?}", e);
        }
    }

    Ok(())
}