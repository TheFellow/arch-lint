# Example (delta)

**Feature-Specific API Access with Pattern Matching:**
- **Pattern-Based Rules**: Uses `{feature}` placeholder for dynamic pattern matching
- **Feature-Specific Access**: API can access feature-specific application modules
- **Self-Reference Allowance**: Features can reference themselves through the pattern

**Enforced Rules:**
- **Feature-Specific Restrictions**: Packages are forbidden from importing `app/{feature}/**` modules
- **API Exception**: Only `api/**` packages can import feature-specific modules
- **Self-Reference Exception**: `app/{feature}/**` modules can import themselves (same feature)

**Pattern Matching Features:**
- **Dynamic Patterns**: `{feature}` acts as a placeholder that matches actual feature names
- **Self-Referential Logic**: The same `{feature}` value in both forbid and except creates self-import allowance

```mermaid
graph TD
    subgraph "Pattern-Based Architecture - example/delta"
        subgraph "API Layer"
            API[bookstore/api/**]
        end

        subgraph "Features"
            F1[bookstore/app/feature1/**]
            F2[bookstore/app/feature2/**]
        end
    end

%% Allowed patterns
    API -->|"✅ Can import<br/>any {feature}"| F1
    API -->|"✅ Can import<br/>any {feature}"| F2
    F1 -->|"✅ Self-reference<br/>same {feature}"| F1
    F2 -->|"✅ Self-reference<br/>same {feature}"| F2

%% Forbidden patterns
    F1 -.->|"❌ Cross-feature<br/>different {feature}"| F2
```