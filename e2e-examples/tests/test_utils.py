from e2e_examples.utils import is_even


def test_is_even():
    for i in range(0, 100, 2):
        assert is_even(i)
    for i in range(1, 100, 2):
        assert not is_even(i)