for i in range(1000):
    with open(f"rac{i}.toml", "w") as f:
        f.write('''[general]
id = "rac"
config_dir = "bench_conf/"

[rac]
ctrl_addr = "127.0.0.1:33333"
addr = "127.0.0.1:PORT"
[[rac.local_algorithms]]
file = "algorithms/module.wasm"
hexhash = "deadbeef"

[log.console]
level = "debug"
                '''.replace("PORT", str(20000 + i)))
