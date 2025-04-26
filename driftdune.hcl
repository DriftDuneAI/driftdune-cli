# driftdune.hcl

driftdune_test_run "example" {
    model           = "gpt-4"
    baseline        = "gpt-3.5"
    prompt_suite    = "suite.json"
    alert_threshold = 0.05
}
