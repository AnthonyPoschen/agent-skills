# Terminal lifecycle

An interactive shell or TUI remains active until its exit action occurs.
Successful startup is not proof of normal termination or cleanup.

## Own the resources

- Prefer a controlled PTY and terminal renderer over new GUI windows for routine automated drives.
- Before launch, record the entry program's lifetime, readiness signal, normal exit action, and deadlines.
- Record each created driver PID, process group, sandbox ID, terminal window ID, and tab ID where applicable.
- Close only resources created by this drive. Keep existing user windows and tabs intact.
- If ownership cannot be established, report the exact unresolved resources instead of broad cleanup.

## Drive normal exit

1. Start the real application through its normal user entry point.
2. Wait for actual guest or application readiness within a bounded deadline.
3. Send a distinctive action through the application's input path.
4. Capture the terminal frame while the driver continues to drain PTY output.
5. Send the documented normal exit action through that same input path.
6. Wait for the entry program's exit status within a bounded deadline.
7. Record whether the application stopped its VM and helpers before driver recovery.
8. Render and inspect the captured frame after the interactive process exits.
9. Restore terminal state and close the drive's GUI tabs or windows.
10. Recover only exact Run-owned resources that the application left behind.
11. Confirm that retained evidence remains available after cleanup.

For a shell-based sandbox, send `exit 0` to the guest shell rather than cancelling the host wrapper.
For a harness TUI, use its documented quit action and verify child-server teardown where applicable.
Do not wait indefinitely for an interactive program to finish by itself.

## Reject false proof

- Separate PTY input echo from application output. Do not accept a marker or field that appears in the command text as proof of its result.
- Keep output markers and machine-readable response labels out of echoed input, or parse a response boundary that excludes input echo.
- Test the proof parser with an echoing PTY that shows the command before the application response.
- Keep draining PTY output while an interactive process runs. A slow screenshot renderer can fill the terminal buffer and block the process before it reads more input.
- Save the live terminal bytes during the drive, then run slow image conversion after the entry process exits.
- User cancellation is an interruption test, not a successful normal-exit test.
- A timeout, missing readiness, or a leaked owned terminal fails the affected criterion.
- Cleanup after failure does not turn the failed drive into a pass.
- Record application cleanup before any harness stop, kill, or remove action. Keep forced recovery evidence separate from product cleanup proof.
- Keep the PTY session leader alive during process-leak checks. Do not let the test harness send SIGHUP that hides an application helper leak.
- Require the golden-path recipe to finish without user intervention.
- Exercise failure and interruption cleanup separately from normal exit.

Protect the driver through an executable check at its lifecycle boundary.
Use a representative interactive process that stays alive until it receives an explicit exit command.
Confirm that a hung readiness or exit path fails instead of reporting success.
Re-drive the real application after a verifier fix before claiming the criterion passes.
