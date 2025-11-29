# Implementation Plan

- [x] 1. Set up project structure and core configuration
  - Create new directory structure (src/edu_agent with subdirectories)
  - Set up pyproject.toml with project metadata and dependencies
  - Create requirements.txt and requirements-dev.txt
  - Implement core configuration module with Pydantic Settings
  - _Requirements: 1.1, 1.5, 5.1, 9.1, 9.3_

- [ ] 1.1 Write property test for configuration defaults
  - **Property 15: Configuration default values**
  - **Validates: Requirements 5.2**

- [ ]* 1.2 Write property test for no hardcoded secrets
  - **Property 16: No hardcoded secrets**
  - **Validates: Requirements 5.3**

- [x] 2. Implement core components (LLM, embeddings, vector store)
  - Create LLM factory with model selection and caching
  - Create embeddings factory
  - Implement VectorStoreManager for multi-tenant vector stores
  - Implement CheckpointerManager for state persistence
  - _Requirements: 2.5, 3.4, 5.5_

- [ ]* 2.1 Write property test for LLM dependency injection
  - **Property 5: LLM dependency injection**
  - **Validates: Requirements 2.5**

- [ ]* 2.2 Write property test for tenant-specific configuration
  - **Property 18: Tenant-specific configuration**
  - **Validates: Requirements 5.5**

- [x] 3. Define state schemas
  - Create base State TypedDict
  - Create CustomerServiceState with all required fields
  - Create SuggestionState
  - Add type annotations and documentation for all fields
  - _Requirements: 4.1, 4.2, 4.4, 4.5_

- [ ]* 3.1 Write property test for state type definitions
  - **Property 10: State type definition**
  - **Validates: Requirements 4.1**

- [ ]* 3.2 Write property test for state field annotations
  - **Property 11: State field type annotations**
  - **Validates: Requirements 4.2**

- [ ]* 3.3 Write property test for optional state fields
  - **Property 13: Optional state fields**
  - **Validates: Requirements 4.4**

- [x] 4. Implement tools layer
  - Create kb_tools.py with retrieve_kb tool using @tool decorator
  - Create order_tools.py with query_order and execute_sql tools
  - Create handoff_tools.py with handoff_to_human tool
  - Implement error handling in all tools
  - Add dependency injection for external resources
  - _Requirements: 3.1, 3.2, 3.3, 3.4_

- [ ]* 4.1 Write property test for tool standardization
  - **Property 6: Tool standardization**
  - **Validates: Requirements 3.1**

- [ ]* 4.2 Write property test for tool dependency injection
  - **Property 7: Tool dependency injection**
  - **Validates: Requirements 3.2**

- [ ]* 4.3 Write property test for tool error handling
  - **Property 8: Tool error handling**
  - **Validates: Requirements 3.3**

- [ ]* 4.4 Write property test for tool configuration usage
  - **Property 9: Tool configuration usage**
  - **Validates: Requirements 3.4**

- [ ] 5. Implement customer service agent nodes
  - Create intent_node with keyword and LLM-based routing
  - Create kb_node with vector search and answer generation
  - Create order_node with SQL generation and NLG
  - Create direct_node for simple responses
  - Create handoff_node for human escalation
  - Ensure all nodes accept State and return dict
  - _Requirements: 2.2_

- [ ]* 5.1 Write property test for node function signatures
  - **Property 2: Node function signature consistency**
  - **Validates: Requirements 2.2**

- [ ] 6. Implement customer service agent routing
  - Create decide_after_intent routing function
  - Create decide_after_kb routing function
  - Ensure routing functions return valid node names
  - _Requirements: 2.3_

- [ ]* 6.1 Write property test for routing function validity
  - **Property 3: Routing function validity**
  - **Validates: Requirements 2.3**

- [ ] 7. Build customer service agent graph
  - Create CustomerServiceAgent class
  - Build StateGraph with all nodes
  - Add conditional edges with routing functions
  - Configure checkpointer for state persistence
  - Compile graph
  - _Requirements: 2.1, 2.4_

- [ ]* 7.1 Write property test for agent StateGraph consistency
  - **Property 1: Agent StateGraph consistency**
  - **Validates: Requirements 2.1**

- [ ]* 7.2 Write property test for checkpointer configuration
  - **Property 4: Checkpointer configuration**
  - **Validates: Requirements 2.4**

- [ ] 8. Implement suggestion agent
  - Create SuggestionAgent with ReAct pattern
  - Implement suggestion generation logic
  - Add tools for KB search and order lookup
  - _Requirements: 2.1, 2.2_

- [ ] 9. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 10. Implement service layer
  - Create ChatService with process_message method
  - Create AuthService with login/register/logout methods
  - Create VectorService with add/delete/search methods
  - Ensure services encapsulate agent invocations
  - _Requirements: 6.2_

- [ ]* 10.1 Write property test for service layer encapsulation
  - **Property 20: Service layer encapsulation**
  - **Validates: Requirements 6.2**

- [ ] 11. Define API models
  - Create ChatRequest and ChatResponse Pydantic models
  - Create LoginRequest, RegisterRequest, AuthResponse models
  - Create VectorAddRequest, VectorDeleteRequest models
  - Create UserRequest, RoleRequest models
  - Add validation rules and documentation
  - _Requirements: 6.3_

- [ ]* 11.1 Write property test for Pydantic request validation
  - **Property 21: Pydantic request validation**
  - **Validates: Requirements 6.3**

- [ ] 12. Implement API dependencies
  - Create get_current_user dependency
  - Create get_current_tenant dependency
  - Create get_chat_service dependency
  - Create require_api_key dependency
  - _Requirements: 6.5_

- [ ]* 12.1 Write property test for authentication dependency injection
  - **Property 23: Authentication dependency injection**
  - **Validates: Requirements 6.5**

- [ ] 13. Implement chat API routes
  - Create /chat endpoint with ChatRequest/ChatResponse
  - Create /suggest/{thread_id} SSE endpoint
  - Create /greet endpoint
  - Ensure routes delegate to ChatService
  - _Requirements: 6.1, 10.1_

- [ ]* 13.1 Write property test for API route separation
  - **Property 19: API route separation**
  - **Validates: Requirements 6.1**

- [ ]* 13.2 Write property test for API endpoint backward compatibility
  - **Property 27: API endpoint backward compatibility**
  - **Validates: Requirements 10.1**

- [ ] 14. Implement auth API routes
  - Create /api/auth/login endpoint
  - Create /api/auth/register endpoint
  - Create /api/auth/logout endpoint
  - Create /api/auth/me endpoint
  - _Requirements: 6.1_

- [ ] 15. Implement admin API routes
  - Create /api/users CRUD endpoints
  - Create /api/roles CRUD endpoints
  - Create /api/v1/vectors/items endpoints
  - Create /models/list and /models/switch endpoints
  - _Requirements: 6.1_

- [ ] 16. Implement middleware
  - Create request ID middleware
  - Create security/redaction middleware
  - Create metrics middleware
  - Create error handling middleware with unified format
  - _Requirements: 6.4_

- [ ]* 16.1 Write property test for unified error response format
  - **Property 22: Unified error response format**
  - **Validates: Requirements 6.4**

- [ ] 17. Implement database access layer
  - Create orders.py with order query functions
  - Create support.py with unanswered questions tracking
  - Create auth.py with user/role management
  - Use parameterized queries to prevent SQL injection
  - _Requirements: 10.3_

- [ ]* 17.1 Write property test for data format backward compatibility
  - **Property 29: Data format backward compatibility**
  - **Validates: Requirements 10.3**

- [ ] 18. Implement utilities
  - Create text processing utilities (clean_input, parse_order_id)
  - Create time utilities
  - Create validators
  - _Requirements: 3.5_

- [ ] 19. Migrate prompts
  - Create prompts/intent.py with INTENT_PROMPT
  - Create prompts/rag.py with RAG_PROMPT_TEMPLATE
  - Create prompts/order.py with ORDER_NLG_PROMPT_TEMPLATE
  - Create prompts/suggestion.py with SUGGEST_QUESTIONS_PROMPT_TEMPLATE
  - _Requirements: 1.1_

- [ ] 20. Create FastAPI application
  - Create main.py with FastAPI app initialization
  - Register all routers
  - Add middleware
  - Configure CORS
  - Add health check endpoint
  - _Requirements: 6.1_

- [ ] 21. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 22. Add environment variable backward compatibility
  - Ensure all original env var names are supported
  - Add aliases for any renamed variables
  - Document any changes in migration guide
  - _Requirements: 10.2_

- [ ]* 22.1 Write property test for environment variable backward compatibility
  - **Property 28: Environment variable backward compatibility**
  - **Validates: Requirements 10.2**

- [ ] 23. Implement feature flags
  - Create feature flag system
  - Add flags for new features
  - Ensure flags default to disabled
  - _Requirements: 10.5_

- [ ]* 23.1 Write property test for feature flag control
  - **Property 30: Feature flag control**
  - **Validates: Requirements 10.5**

- [ ] 24. Create documentation
  - Write README.md with installation and usage instructions
  - Create docs/architecture.md explaining design decisions
  - Create docs/api.md with API documentation
  - Create docs/development.md with development guidelines
  - Create docs/migration.md with migration guide
  - Add module and function docstrings
  - _Requirements: 8.1, 8.2, 8.3, 8.4, 10.4_

- [ ]* 24.1 Write property test for module docstrings
  - **Property 24: Module docstrings**
  - **Validates: Requirements 8.2**

- [ ]* 24.2 Write property test for function parameter documentation
  - **Property 25: Function parameter documentation**
  - **Validates: Requirements 8.3**

- [ ] 25. Create examples
  - Create examples/basic_chat.py demonstrating simple chat
  - Create examples/custom_agent.py showing how to extend agents
  - Create examples/multi_tenant.py showing tenant configuration
  - _Requirements: 8.5_

- [ ] 26. Set up testing infrastructure
  - Create tests/conftest.py with shared fixtures
  - Create tests/fixtures/factories.py with test data generators
  - Set up pytest configuration
  - Set up Hypothesis for property-based testing
  - _Requirements: 7.1, 7.5_

- [ ]* 26.1 Write unit tests for tools
  - Test each tool with mocked dependencies
  - Test error handling paths
  - _Requirements: 7.2_

- [ ]* 26.2 Write integration tests for agents
  - Test complete agent workflows
  - Test state transitions
  - _Requirements: 7.3_

- [ ]* 26.3 Write API tests
  - Test all endpoints with TestClient
  - Test authentication flows
  - Test error scenarios
  - _Requirements: 7.4_

- [ ] 27. Add type checking
  - Configure mypy for static type checking
  - Add mypy to CI/CD pipeline
  - Fix any type errors
  - _Requirements: 4.3_

- [ ]* 27.1 Write property test for state type consistency
  - **Property 12: State type consistency**
  - **Validates: Requirements 4.3**

- [ ] 28. Add dependency version constraints
  - Pin major dependency versions in requirements.txt
  - Document version constraints rationale
  - _Requirements: 9.2_

- [ ]* 28.1 Write property test for dependency version constraints
  - **Property 26: Dependency version constraints**
  - **Validates: Requirements 9.2**

- [ ] 29. Final integration testing
  - Run full test suite
  - Test with real LLM and vector store
  - Test multi-tenant scenarios
  - Verify backward compatibility with existing deployments
  - _Requirements: 10.1, 10.2, 10.3_

- [ ] 30. Final checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.
