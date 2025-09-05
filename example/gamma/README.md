# Example (gamma)

**Feature Isolation with Exemptions:**
- **Feature Separation**: Features should be isolated and cannot import each other
- **Selective Exemptions**: Specific features can be granted exemptions from isolation rules
- **Module Independence**: Promotes loose coupling between feature modules

**Enforced Rules:**
- **No Cross-Feature Imports**: Features under `feature/*` are forbidden from importing other features (`feature/*`)
- **Feature B Exemption**: Feature B is exempt from the isolation rule and can be imported by others

**Current Violations Detected by arch-lint:**
1. **Feature Cross-Import**: `feature/A/feat-A.go` imports `feature/B`, but Feature B is exempted so this should be allowed

```mermaid
graph TD
    subgraph "Feature Isolation - example/gamma"
        subgraph "Feature A"
            FA[feature/A/feat-A.go]
        end

        subgraph "Feature B (Exempt)"
            FB[feature/B/feat-b.go]
        end
    end

%% Current imports (Feature B is exempt, so this should be allowed)
    FA -->|"✅ Allowed<br/>(Feature B is exempt)"| FB
```