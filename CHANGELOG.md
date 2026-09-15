# Changelog

## [0.13.0](https://github.com/RomanAgaltsev/undergo/compare/v0.12.0...v0.13.0) (2026-09-15)


### Features

* **tasks:** M11 — the floor sweep, every track at five ([#17](https://github.com/RomanAgaltsev/undergo/issues/17)) ([e79a6cc](https://github.com/RomanAgaltsev/undergo/commit/e79a6cc83a6819df9064b0ac5c80a451c7c24932))

## [0.12.0](https://github.com/RomanAgaltsev/undergo/compare/v0.11.0...v0.12.0) (2026-09-15)


### Features

* **tasks:** M10 — sched and memmodel tracks ([#15](https://github.com/RomanAgaltsev/undergo/issues/15)) ([a635e5d](https://github.com/RomanAgaltsev/undergo/commit/a635e5d864d5f345a727f6d95121e13771f20d0f))

## [0.11.0](https://github.com/RomanAgaltsev/undergo/compare/v0.10.0...v0.11.0) (2026-09-14)


### Features

* **tasks:** the machine — iface, compiler, asm, edges ([#13](https://github.com/RomanAgaltsev/undergo/issues/13)) ([1c4f40b](https://github.com/RomanAgaltsev/undergo/commit/1c4f40b7a95966ea2fffea07cd7b447b8f4ffec2))

## [0.10.0](https://github.com/RomanAgaltsev/undergo/compare/v0.9.0...v0.10.0) (2026-09-14)


### Features

* **tasks:** lifetime and collection — weak and gc ([#11](https://github.com/RomanAgaltsev/undergo/issues/11)) ([4b3ce0e](https://github.com/RomanAgaltsev/undergo/commit/4b3ce0ebde1ab69e3ea9ec641761abb89b14be58))

## [0.9.0](https://github.com/RomanAgaltsev/undergo/compare/v0.8.0...v0.9.0) (2026-09-13)


### Features

* **ci:** race the machinery and every reference solution ([#9](https://github.com/RomanAgaltsev/undergo/issues/9)) ([187ee6a](https://github.com/RomanAgaltsev/undergo/commit/187ee6ae81c26019e2c3be2d122c5018259b79cf))

## [0.8.0](https://github.com/RomanAgaltsev/undergo/compare/v0.7.0...v0.8.0) (2026-09-13)


### Features

* **tasks:** the language surface — generics, iter, reflect ([#7](https://github.com/RomanAgaltsev/undergo/issues/7)) ([d22b17e](https://github.com/RomanAgaltsev/undergo/commit/d22b17ece299ef394d1beac829b057518ddb4316))

## [0.7.0](https://github.com/RomanAgaltsev/undergo/compare/v0.6.0...v0.7.0) (2026-09-13)


### Features

* **tasks:** alloc/02-allocs-per-op ([b21837f](https://github.com/RomanAgaltsev/undergo/commit/b21837fcf3f18b4bedff98ebbcb12fa1367a186b))
* **tasks:** alloc/03-what-escapes ([68db8c7](https://github.com/RomanAgaltsev/undergo/commit/68db8c7af1175dfb2581a24160c8271ab1ea7bb2))
* **tasks:** alloc/04-boxing-cache ([1f3ba94](https://github.com/RomanAgaltsev/undergo/commit/1f3ba9427eaebee1e416f038549362ab5f718eeb))
* **tasks:** layout/02-reorder-to-shrink ([3b5e7b5](https://github.com/RomanAgaltsev/undergo/commit/3b5e7b59382533c6a8a0dcccc01d82a4af23c0a2))
* **tasks:** layout/03-false-sharing ([37e7e05](https://github.com/RomanAgaltsev/undergo/commit/37e7e0570dccf20c2ff7aac6806fd6bd640806e6))
* **tasks:** types/01-cap-growth ([f84c2f1](https://github.com/RomanAgaltsev/undergo/commit/f84c2f199c9dff8e0918c4ee5f22662f4e124c09))
* **tasks:** types/02-aliasing ([7635463](https://github.com/RomanAgaltsev/undergo/commit/76354638a436a5cd08927ea8995aa6904c8e8c9f))
* **tasks:** types/03-zero-copy-strings ([8e783af](https://github.com/RomanAgaltsev/undergo/commit/8e783aff1d1c6d5f10393b2a86ecda85c4d5ab50))
* **tasks:** types/04-map-memory ([1f6b148](https://github.com/RomanAgaltsev/undergo/commit/1f6b148aa5bc254740ab6d98fe2adfc39a0c68e0))


### Documentation

* the memory tracks are complete ([9e626ea](https://github.com/RomanAgaltsev/undergo/commit/9e626eaebdf75ab0840ebcad84449e53c592c1e8))

## [0.6.0](https://github.com/RomanAgaltsev/undergo/compare/v0.5.1...v0.6.0) (2026-09-11)


### ⚠ BREAKING CHANGES

* all 135 review task ids changed. Any progress record naming an old id no longer resolves to a task.

### Features

* nest review drills under tasks/review/ by category ([608696a](https://github.com/RomanAgaltsev/undergo/commit/608696ad34d32d20d659d26d03b4bf791432c784))

## [0.5.1](https://github.com/RomanAgaltsev/undergo/compare/v0.5.0...v0.5.1) (2026-09-11)


### Bug Fixes

* **tasks:** layout/01-struct-padding attributes field order to gc, not to Go ([fc62330](https://github.com/RomanAgaltsev/undergo/commit/fc62330d91953d5ef2972a2912edf3e9b6030d1b))

## [0.5.0](https://github.com/RomanAgaltsev/undergo/compare/v0.4.0...v0.5.0) (2026-09-11)


### Features

* **radar:** a checked format for the release radar ([fb7a223](https://github.com/RomanAgaltsev/undergo/commit/fb7a223b047cb6962e6f3810808b9f42b17d47c6))


### Documentation

* pool the E1 radar findings into the candidate backlog ([642d9f7](https://github.com/RomanAgaltsev/undergo/commit/642d9f77f2012be51a1ce5ffcd8eb656d2901275))
* **radar:** triage Go 1.23 ([ddcd73b](https://github.com/RomanAgaltsev/undergo/commit/ddcd73b3997fc0a0a12f226d4b03b959b1b90751))
* **radar:** triage Go 1.24 ([f606895](https://github.com/RomanAgaltsev/undergo/commit/f606895662c434d7f0fe3d810a0bac5de9872695))
* **radar:** triage Go 1.25 ([f6eca31](https://github.com/RomanAgaltsev/undergo/commit/f6eca31d7335c8745434b607d952626c13c8971d))
* **radar:** triage Go 1.26 ([4f48402](https://github.com/RomanAgaltsev/undergo/commit/4f4840214b52193caa72e4015795340625803b78))
* **radar:** triage Go 1.27 ([78c0b05](https://github.com/RomanAgaltsev/undergo/commit/78c0b05b7da27795964c81a7ac1cd039506aa69a))
