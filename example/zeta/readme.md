# Example (zeta)

**Proper Clean Architecture with Utilities:**
- **Utilities** can be imported by Controllers, Use Cases, and Infrastructure
- **Domain** remains completely independent (no imports allowed)
- **Controllers** should depend on Use Cases, not Infrastructure directly
- **Use Cases** orchestrate domain logic and can use utilities
- **Infrastructure** implements domain interfaces and can use utilities

**Enforced Rules:**
- **Domain Independence**: Domain layer (`domain/**`) is forbidden from importing anything (`"**"`), including utilities
- **Controller Isolation**: Controllers (`controllers/**`) are forbidden from importing infrastructure (`infrastructure/**`)

**Current Violations Detected by arch-lint:**
1. **Controllers → Infrastructure**: `controllers/controller.go` directly imports `infrastructure/db`, bypassing the use case layer
2. **Domain → Util**: `domain/entity.go` imports `util`, violating domain independence


```mermaid
graph TD
%% Clean Architecture Layers
    subgraph "Clean Architecture - example/zeta"
        subgraph "Controllers Layer"
            C[controllers/controller.go]
        end

        subgraph "Use Case Layer"
            UC[usecase/service.go]
        end

        subgraph "Domain Layer"
            D[domain/entity.go]
        end

        subgraph "Infrastructure Layer"
            I[infrastructure/db/repo.go]
        end

        subgraph "Utility Layer"
            U[util/utils.go]
        end
    end

%% Correct Clean Architecture Flow
    C -->|✅ Should import| UC
    UC -->|✅ Can import| D
    UC -->|✅ Can import| I
    I -->|✅ Can import| D

%% Current Violations (what arch-lint catches)
    C -.->|❌ VIOLATION<br/>controllers → infrastructure| I
    D -.->|❌ VIOLATION<br/>domain → util| U


```