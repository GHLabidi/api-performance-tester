import sys
from analyzer import Analyzer



# main
if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python3 generate_report.py <unique_test_name>")
        sys.exit(1)
    # create analyzer
    analyzer = Analyzer(sys.argv[1])
    # load data
    analyzer.load_data()
    # analyze requests per second
    test_description = """
                        This represents the number of requests processed by the server per second.
                        """
    analyzer.analyze_requests_per_second(test_description=test_description)
    # analyze lookup duration
    test_description = """
                        This represents the lookup duration of the client's ip in the server.
                        In this test, the server is storing the ip addresses in a hash table along with the number of requests made by this ip.
                        """
                        
    analyzer.analyze_latency('QueryDuration', 'Query Duration', test_description=test_description)
    # analyze request duration
    test_description = """
                        This represents the whole request duration from the moment the request is received by the server until the response is sent back to the client.
                        """
    analyzer.analyze_latency('RequestDuration', 'Request Duration', test_description=test_description)
    # create report
    analyzer.create_test_report_html()
