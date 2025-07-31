The subfolders correspond to the different application modules and each contains the test cases relevant to that module.

The SDK was designed so that all the internal logic is fully encapsulated and abstracted away. Because of this, there isn’t a way to test the user-facing interfaces directly—only the internal logic can be tested. This is mainly because the AWS clients are created inside the SDK itself and aren’t injected or passed in via interfaces, which makes mocking or substituting them in tests difficult. Internally custom fetchers are passed via interfaces to the actual matrics functions that calculate the metrics and thus they can be tested.

As a result, the user interface parts are currently only tested manually.
