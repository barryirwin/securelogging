# Benchmarks

This is to keep track of how slow/fast the server becomes as more features are added.

## Specs

Done in a laptop with the following specs:

- Intel(R) Core(TM) i7-6820HQ CPU @ 2.70GHz
- 16 GB RAM
- 512 GB SSD

With the following settings in the config file:

```bash
# Perfomance tweaking
# NOTE: Increasing buffer sizes will increase memory consumption accordingly
# Network packet buffer, increase if the size of the packets are bigger
PacketBufferSize = 2048
# Batch size to store logs, adjust according to number of logs received
FileBatchSize = 1000
# Buffer size for the channel receiving packets (number of packets allowed to be queued), increase if the data ingestion is lagging
WireChannelBuffer = 1
# Buffer size for all processing channels (after they are received), increase if the data writting/influx/decryption is lagging
ProcessingChannelBuffer = 1
```

The laptop then decided to die... so got a Macbook Air M1 with 16GB Ram, which basically runs everything twice as fast.

## Script

There is benchmark script on the `run/` folder that sends 10 000 UDP and 10 000 TCP packets (one and one) with a relatively long string

## Perfomance timeline

As features are added in time, a new category is added here

### Barebones

Just receiving the packets, going through the different channels, but the data is passed as is and storing them to a file, takes around 1s to process all 20 000 packets.

```bash
franco-> run $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m0.952s
user    0m0.397s
sys     0m0.547s
franco-> run $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m0.995s
user    0m0.414s
sys     0m0.574s
```

### Password encryption

After adding the password decryption before storing, the test string was updated with the same as used on the unit test for the decryption, running that string with the barebones version shields same perfomance compared to the older gibberish.

Interestingly, the password encryption does not seem to affect the data ingestion.

```bash
franco-> run $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m0.903s
user    0m0.392s
sys     0m0.502s
franco-> run $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.088s
user    0m0.448s
sys     0m0.625s
```

### Storing data - Txt + Password Encryption

At this stage, the data is parsed to a json structure, then encoded as base64 string, then password encrypted and stored. This added some processing time, but is still within a good digestion rate:

```bash
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.061s
user    0m0.420s
sys     0m0.632s
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.025s
user    0m0.369s
sys     0m0.646s
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.086s
user    0m0.422s
sys     0m0.647s
```

### Storing data - Added compression

Some extra sanity checks were added, but the main feature is the gzip compression of the blob data being stored to disk. This saves around 25% of storage (measured in byte size), but adds processing time. Probably the space saving increases as more "real data" is used later in the future.

```bash
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.182s
user    0m0.515s
sys     0m0.638s
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.211s
user    0m0.436s
sys     0m0.717s
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.146s
user    0m0.505s
sys     0m0.627s
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.170s
user    0m0.455s
sys     0m0.700s
```

Disabling compression with the same code shields:

```bash
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.030s
user    0m0.405s
sys     0m0.616s
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.081s
user    0m0.431s
sys     0m0.634s
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.020s
user    0m0.382s
sys     0m0.629s
```

### Adding SHA256 hash of the previous line to struct

This was added to have a way to reference the previous log in the file, so when the replay mode/decoders is in place, it will be possible to flag if a log entry was removed.

This doesn't seem to add any significant processing time, looks much the same:

```bash
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.052s
user    0m0.405s
sys     0m0.628s
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.058s
user    0m0.488s
sys     0m0.560s
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m0.981s
user    0m0.385s
sys     0m0.587s
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.025s
user    0m0.421s
sys     0m0.596s
```

### Adding RSA encryption

After adding RSA encryption/decription with a long key (8096/1024), to be able to encrypt relatively long messages (up to 800 chars or so), but is not going to work as even a small encrypted and enriched packet will be longer than that.

The code to generate the keys, read, encrypt, etc is on place, but will not be used for now. As reference, a single short string takes (from unit test):

```bash
Encryption took 0.000756258
Decryption took 0.065599442
```

After some back and forth, decided that another approach is needed, AES-256 is the recommended algorithm to use on the NSA, that is good enough encryption for this project. The remaining problem is that now all the packet structure is different (AES256 blob + RSA encrypted AES password), so refactoring the main logic flow and related functions was required.

This increased the string length received significantly, making the benchmark script to overload the network stack somehow, after some tests, UDP seemed to handle the load better, but increasing the number of packets to 100k made the server crash.

Need to investigate this bottlenck later on. Could be also my laptop which has the TPM defective and since I was lucky to get it to start, haven't turned it off or restarted in a month, need to test somewhere else.

Testing with different string sizes, looks like removing the comma that separated delimiter helped somewhat as less load is put on the program, packets are stored with less processing, still 10k packets crash the server. This is probably related to the extra processing of having the packet with this format: decryption or I/O, needs to be investigated:

```bash
# The server can handle 1k TCP / 1k UDP of this packet string with a FileBatchSize = 1000, FileBatchSize = 100 crashes it.
# Increasing this value in the test script to 1500 makes the server crash
pass="PoUu+c2fFLdCZWU/Tu+urQsvQAxgOzfjW1h8Gidf6JjuCfSF7TNcWN0ypQqe4sfGqgKtGFqIe72Ml6V6bdrV4HbIpav0RAXZCWtGAcPW4aJMIUpd69asUzL6QLmrKOlJS3DAol/frJtNg+G6udT9KN6Fd4SBTbqmropxXbim41aZgsQYEKEMVcRAPXczZJm6y7jxJKEATviIfLOSKDmWxF6S7AzaB69bcWn+QJNKfd6nqBcXlgKESwwWjLs48/hWTIohT6Rov5Bzjm+Chyr5fePMUQLoPDKRJqQ1MTmIQmUcP6onFeMVp6uKurnjL02JXiK5dux7YW/tnY7/AqhD1Q=="
# "This is a test" AES-256 encrypted with "my-fancy-password" SHA-256
blob="9a959a392810308a4df45613e82548a375042c6c1d8bacc76b1f89e4fb168d7611eadd7c7ed210e32123"

message="$blob,$pass"
```

Interestingly enough, the benchmark script suggests that the server should be able to handle the load:

```bash
franco->run\ $ time ./stress-test-gibberish.sh
Sent 1000 packets on UDP and 1000 packets on TCP

real    0m0.126s
user    0m0.081s
sys     0m0.044s
franco->run\ $ time ./stress-test-gibberish.sh
Sent 1000 packets on UDP and 1000 packets on TCP

real    0m0.155s
user    0m0.060s
sys     0m0.060s
franco->run\ $ time ./stress-test-gibberish.sh
Sent 1000 packets on UDP and 1000 packets on TCP

real    0m0.125s
user    0m0.062s
sys     0m0.062s
```

This malformed script however points further into the decryption as being the culprit, there is a perfomance loss compared to earlier versions, likely due the string length or shared key?

```bash
# The server can handle 10k TCP / 10k UDP of this packet string with a FileBatchSize = 1000, FileBatchSize = 100 crashes it.
message="9a959a392810308a4df45613e82548a375042c6c1d8bacc76b1f89e4fb168d7611eadd7c7ed210e32123PoUu+c2fFLdCZWU/Tu+urQsvQAxgOzfjW1h8Gidf6JjuCfSF7TNcWN0ypQqe4sfGqgKtGFqIe72Ml6V6bdrV4HbIpav0RAXZCWtGAcPW4aJMIUpd69asUzL6QLmrKOlJS3DAol/frJtNg+G6udT9KN6Fd4SBTbqmropxXbim41aZgsQYEKEMVcRAPXczZJm6y7jxJKEATviIfLOSKDmWxF6S7AzaB69bcWn+QJNKfd6nqBcXlgKESwwWjLs48/hWTIohT6Rov5Bzjm+Chyr5fePMUQLoPDKRJqQ1MTmIQmUcP6onFeMVp6uKurnjL02JXiK5dux7YW/tnY7/AqhD1Q=="
```

```bash
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.485s
user    0m0.648s
sys     0m0.714s
franco->run\ $ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP

real    0m1.406s
user    0m0.668s
sys     0m0.703s
```

Test again on the M1...

```bash
# This is with garbage
~/study/noroff-fdp/run on  main! ⌚ 11:05:31
$ time ./stress-test-gibberish.sh
Sent 10000 packets on UDP and 10000 packets on TCP
./stress-test-gibberish.sh  0.37s user 0.55s system 46% cpu 1.998 total
```

## Replay data

From this part onwards the development was done after the laptop change

```bash
~/study/noroff-fdp/build on  main! ⌚ 20:22:06
$ time ./slog-server-darwin-arm64.bin -read-file slog-data/slog-logging.slog -password my-fancy-password -priv-key keys/storage/slog-storage_rsa
2021-01-15T20:22:07.538+0100    INFO    slogserver/main.go:75   Will try to re-process the file: slog-data/slog-logging.slog
2021-01-15T20:22:07.539+0100    INFO    slogserver/main.go:80   Using the private key at keys/storage/slog-storage_rsa and the password my-fancy-password
2021-01-15T20:22:07.539+0100    INFO    slogserver/main.go:82   And will try to write the output to slog-reprocessed.txt
2021-01-15T20:22:07.539+0100    INFO    server/replayProcessor.go:43    Bear in mind that if the segment read is was not the start of the file, the order of the first log in the file can't be verified. So will always be flagged as suspicious.
2021-01-15T20:22:07.539+0100    INFO    server/replayProcessor.go:44    To solve this, just start slightly before the segment of interest and disregard the first log failure. Or replay the whole original file.
2021-01-15T20:22:07.540+0100    INFO    server/replayProcessor.go:53    Looks like this is the start of the file
2021-01-15T20:22:14.926+0100    INFO    server/replayProcessor.go:99    Replay Summary: {"Total lines processed": 7000, "Replay Suspicious Lines": 0, "Suspicious packets (When data arrived)": 0}
./slog-server-darwin-arm64.bin -read-file slog-data/slog-logging.slog      7.14s user 0.16s system 98% cpu 7.410 total
````
