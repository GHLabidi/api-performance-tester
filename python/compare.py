import sys
from analyzer import Analyzer

def main():
    if len(sys.argv) != 3:
        print("Usage: python3 compare_tests.py <test_1_name> <test_2_name>")
        sys.exit(1)
    # create analyzer

    analyzer = Analyzer(sys.argv[1], sys.argv[2])
    analyzer.compare_tests(['RequestsPerSecond', 'QueryDuration', 'RequestDuration'])

if __name__ == "__main__":
    main()