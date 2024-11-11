# Trinity App Architecture Diagram

Below is a professional architecture diagram of the Trinity App, illustrating the interaction between various components within the system.

```mermaid
graph TD
    %% Clients
    subgraph Client
        User[User / Client]
    end

    %% Application Server
    subgraph Trinity_App
        API[Trinity API Server]
        
        %% Internal Services
        subgraph Internal_Services
            Campaign_Service[Campaign Service]
            Voucher_Service[Voucher Service]
            Purchase_Service[Purchase Service]
        end
        
        %% Common Modules
        subgraph Common_Modules
            Localization_Module[Localization Module]
            Logger[Logger]
            Error_Response[Error Response]
        end
        
        %% API Documentation
        Swagger[Swagger UI]
    end

    %% Database
    subgraph Database
        MongoDB[MongoDB]
    end

    %% Relationships
    User -->|HTTP Requests| API
    API --> Campaign_Service
    API --> Voucher_Service
    API --> Purchase_Service
    API --> Swagger
    API --> Localization_Module
    API --> Logger
    Campaign_Service --> MongoDB
    Voucher_Service --> MongoDB
    Purchase_Service --> MongoDB
    Error_Response --> API
