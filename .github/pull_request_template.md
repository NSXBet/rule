<!-- 
🚀 NSXBet/rule Contribution Guidelines

Thank you for contributing! This template helps ensure high-quality contributions.
Please fill out the sections below and remove any that don't apply.

📋 REQUIREMENTS CHECKLIST (check before submitting):
- [ ] Tests pass (`make test`)
- [ ] Linter passes 100% clean (`make lint`) 
- [ ] Fuzz tests pass (`make fuzz`) for edge case validation
- [ ] Zero allocations maintained in hot paths (if touching core evaluation)
- [ ] Benchmarks included for performance changes
- [ ] Documentation updated if adding features

🔧 DEVELOPMENT COMMANDS:
- `make test` - Run all tests
- `make lint` - Run linter (must be 100% clean)
- `make bench` - Run benchmarks
- `make fuzz` - Run fuzz tests for edge case detection
- `make format` - Format code

⚡ PERFORMANCE REQUIREMENTS:
- Core evaluation must remain under 100ns
- Zero allocations during rule evaluation (0 allocs/op)
- Thread-safe for concurrent usage

📚 CONTRIBUTION TYPES:
- 🐛 Bug fixes: Include test case reproducing the issue
- ⚡ Performance: Include before/after benchmarks
- 🔧 Features: Include tests, benchmarks, and documentation
- 📖 Docs: Ensure accuracy and clarity

🧪 TESTING:
- Add tests for new features
- Include edge cases and error scenarios
- Property-to-property comparisons should have comprehensive coverage
- DateTime operations should include timezone edge cases
- Run fuzz tests (`make fuzz`) for comprehensive edge case detection

📖 DOCUMENTATION:
- Update README if adding features
- Include code examples for new operators
- Update compatibility section if changing behavior
- Follow existing documentation style

🔒 SECURITY:
- Never commit secrets or keys
- Follow defensive programming practices
- Handle edge cases gracefully (don't panic)

For detailed guidelines, see the Contributing section in README.md
-->

## 📝 Summary

<!-- Add here a brief description of what this PR does.--->

## ⚡ Performance Impact

<!-- For performance-related changes, include benchmarks -->

- [ ] No performance impact
- [ ] Performance improved (include benchmarks)
- [ ] Performance regression (justify why acceptable)

## ✅ Checklist

- [ ] Code follows project style
- [ ] Self-reviewed the code
- [ ] Tests added for new functionality
- [ ] README updated (if required)
