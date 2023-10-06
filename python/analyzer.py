import json
import sys
import pandas as pd
import numpy as np
import plotly.graph_objs as go
from plotly.subplots import make_subplots


class Analyzer:

    # data_path: path to the data file
    data_path = None
    test_1_name = None
    test_2_name = None
    server_url = None
    test_start_time = None
    title = ''
    description = ''
    concurrent_requests = None
    test_duration = None
    total_requests = None
    successful_requests = None
    failed_requests = None
    requests_per_second = None
    test_mode = None
    df = None
    meta = None
    graphs = []
    
    # constructor for generating report for one test or comparing two tests
    def __init__(self, test_1_name = None, test_2_name = None):
        if test_1_name is None and test_2_name is None:
            print("Usage: python3 generate_report.py <unique_test_name> or python3 compare_tests.py <test_1_name> <test_2_name>")
            sys.exit(1)
        # check if the user wants to compare two tests
        if test_1_name is not None and test_2_name is not None:
            self.test_1_name = test_1_name
            self.test_2_name = test_2_name
            self.test_1_data_path = 'data/' + test_1_name + '/'
            self.test_2_data_path = 'data/' + test_2_name + '/'
            
        # else the user wants to generate report for one test
        if test_1_name is not None:
            test_unique_name = test_1_name
        self.data_path = 'data/' + test_unique_name + '/'
       
        # read meta data
        # print("Loading meta data from: " + self.data_path + 'benchmark.json')
        meta = json.load(open(self.data_path + 'benchmark.json'))
        self.server_url = meta['request_url']
        # get test start time and convert it to datetime
        self.test_start_time = pd.to_datetime(meta['test_start_time'], unit='ns')
        # convert to readable format YYYY-MM-DD HH:MM:SS
        self.test_start_time = self.test_start_time.strftime("%Y-%m-%d %H:%M:%S")
        # extract test information
        self.title = meta['test_display_name']
        self.description = meta['test_description']
        self.test_mode = meta['test_mode']
        self.concurrent_requests = meta['concurrent_requests']
        self.test_duration = meta['test_duration']
        self.total_requests = meta['total_requests']
        self.successful_requests = meta['successful_requests']
        self.failed_requests = meta['failed_requests']
        self.requests_per_second = meta['requests_per_second']
        
        
    def load_data(self):
        # print("Loading data from: " + self.data_path + 'data.csv')

        if self.total_requests is None or self.total_requests == 0:
            raise ValueError("Could not load data. Total requests is 0")
        # types
        data_types = {
            'StartTime': np.int64,  
            'EndTime': np.int64, 
            'QueryDuration': np.int64,
            'RequestDuration': np.int64,
        }
        self.df = pd.read_csv(self.data_path+'data.csv', header=None, names=data_types.keys(), dtype=data_types)
        # setting datetime fields
        self.df['StartTime'] = pd.to_datetime(self.df['StartTime'], unit='ns')
        self.df['EndTime'] = pd.to_datetime(self.df['EndTime'], unit='ns')
        
    def create_fig(self, chart_data):
        fig = make_subplots(rows=1, cols=2, subplot_titles=("Mean " + chart_data["display_name"] +  " Per Second", "Summary Information"))

        # left figure
        fig.add_trace(go.Scatter(x=chart_data["average_data_per_second_chart_data"]["x"],
                                y=chart_data["average_data_per_second_chart_data"]["y"],
                                mode=chart_data["average_data_per_second_chart_data"]["mode"],
                                name=chart_data["average_data_per_second_chart_data"]["name"],),
                    row=1, col=1)
        fig.update_xaxes(title_text=chart_data["average_data_per_second_chart_data"]["xaxis_title_text"] , row=1, col=1)
        fig.update_yaxes(title_text=chart_data["average_data_per_second_chart_data"]["yaxis_title_text"], row=1, col=1)

        # right figure
        fig.add_trace(go.Bar(
                        x=chart_data["data_description_chart_data"]["x"],
                        y=chart_data["data_description_chart_data"]["y"],
                        text=chart_data["data_description_chart_data"]["text"],
                        textposition='outside',
                        #marker_color=chart_data["data_description_chart_data"]["marker_color"],
                        name=chart_data["data_description_chart_data"]["name"],),
                    row=1, col=2)
        fig.update_xaxes(title_text=chart_data["data_description_chart_data"]["xaxis_title_text"], row=1, col=2)
        fig.update_yaxes(title_text=chart_data["data_description_chart_data"]["yaxis_title_text"], row=1, col=2)


        # setup the layout for the subplots
        fig.update_layout(title_text=chart_data["display_name"] + " Analysis", showlegend=False)
        fig.update_layout(height=600, width=1200)
        return fig
    
    def analyze_latency(self, field_name, display_name, test_description = ''):
        if self.df is None:
            raise ValueError("The data was not loaded")
        
        
        # Analyze Mean Latency Per Second
        ## copy df and aggregate 'field_name' per second
        df_resampled = self.df.copy()
        df_resampled.set_index('StartTime', inplace=True)
        df_resampled = df_resampled.resample('1S').agg({field_name: 'mean'})
        ## reset index to show seconds passed instead of StartTime which is datetime
        start_time = df_resampled.index.min()
        df_resampled['SecondsPassed'] = (df_resampled.index - start_time).total_seconds()
        df_resampled.set_index('SecondsPassed', inplace=True)
        
    
        # convert nanoseconds to milliseconds for better readability
        df_resampled[field_name] = df_resampled[field_name] / 1000000
        
        
        # Analyze Summary Information
        ## copy df and extract 'field_name' column
        df_copy = self.df[field_name]
        ## get data description
        data_description = df_copy.describe(percentiles=[.25, .5, .75, .90, .95, .99]).astype('int64')
        # print("Data description")
        # print(data_description)
        ## extract analysis    
        min = data_description['min']
        min = min / 1000000 # convert to milliseconds

        mean = data_description['mean']
        mean = mean / 1000000 # convert to milliseconds

        std = data_description['std']
        std = std / 1000000 # convert to milliseconds

        max = data_description['max']
        max = max / 1000000 # convert to milliseconds

        p25 = data_description['25%']
        p25 = p25 / 1000000 # convert to milliseconds

        p50 = data_description['50%']
        p50 = p50 / 1000000 # convert to milliseconds

        p75 = data_description['75%']
        p75 = p75 / 1000000 # convert to milliseconds

        p90 = data_description['90%']
        p90 = p90 / 1000000

        p95 = data_description['95%']
        p95 = p95 / 1000000

        p99 = data_description['99%']
        p99 = p99 / 1000000
    
        ## html summary
        html_div = f"""
                    <div class='card rounded-xl m-10 p-10 border-2'>
                        <p class="text-lg font-bold italic">Summary Information</p>
                        <p class="italic">Mean Latency: <b>{mean}</b> milliseconds</p>
                        <p class="italic">Standard Deviation: <b>{std}</b> milliseconds</p>
                        <p class="italic">25th Percentile: <b>{p25}</b> milliseconds</p>
                        <p class="italic">50th Percentile: <b>{p50}</b> milliseconds</p>
                        <p class="italic">75th Percentile: <b>{p75}</b> milliseconds</p>
                        <p class="italic">90th Percentile: <b>{p90}</b> milliseconds</p>
                        <p class="italic">95th Percentile: <b>{p95}</b> milliseconds</p>
                        <p class="italic">99th Percentile: <b>{p99}</b> milliseconds</p>
                        <p class="italic">Request with lowest latency: <b>{min}</b> milliseconds</p>
                        <p class="italic">Request with highest latency: <b>{max}</b> milliseconds</p>
                    </div>
                    """
        analyzed_data = {
            'name' : field_name, # used to identify the file name to be saved
            'display_name': display_name, 
            'data_description_chart_data' : {
                'x': ['25%', '50%', '75%', '90%', '95%', '99%'],
                'y': [p25, p50, p75 , p90, p95, p99],
                'text': [p25, p50, p75 , p90, p95, p99], 
                'textposition': 'outside',
                #'marker_color': ['blue', 'yellow', 'red' , 'green', 'orange', 'purple'],
                'name': 'Data Description',
                'xaxis_title_text': 'Percentile',
                'yaxis_title_text': display_name + '(ms)',
            },
        
            
            'average_data_per_second_chart_data': {
                'x': df_resampled.index.to_list(),
                'y': df_resampled[field_name].to_list(),
                'mode': 'lines',
                'name': 'Mean ' + display_name,
                'xaxis_title_text': 'Time (seconds)',
                'yaxis_title_text': display_name + '(ms)',

            }
        }

        # left figure
        fig = self.create_fig(analyzed_data)
        self.graphs.append({
            'title': display_name + " Analysis",
            'description': test_description,
            'html_summary': html_div,
            'fig': fig,
            'analyzed_data': analyzed_data
        })

    def analyze_requests_per_second(self, test_description = ''):
        # copy and aggregate per second
        df_copy = self.df.copy()
        df_copy.set_index('StartTime', inplace=True)
        df_copy = df_copy.resample('1S').count() 
        # reset index to show seconds passed instead of StartTime which is datetime
        start_time = df_copy.index.min()
        df_copy['SecondsPassed'] = (df_copy.index - start_time).total_seconds()
        df_copy.set_index('SecondsPassed', inplace=True)
        # rename columns and drop unnecessary ones
        df_copy.rename(columns={'EndTime': 'RequestsPerSecond'}, inplace=True)
        df_copy.drop(columns=['QueryDuration', 'RequestDuration'], inplace=True)
        
        data_description = df_copy.describe(percentiles=[.25, .5, .75, .90, .95, .99]).astype('int64')
        # print("Data description")
        # print(data_description)
        data_description_dict = data_description.to_dict()
        # print("Data description dict")
        # print(data_description_dict)
        # extract x and y for percentile chart
        x = ['25%', '50%', '75%', '90%', '95%', '99%']
        y = [data_description_dict['RequestsPerSecond']['25%'], data_description_dict['RequestsPerSecond']['50%'], data_description_dict['RequestsPerSecond']['75%'] , data_description_dict['RequestsPerSecond']['90%'], data_description_dict['RequestsPerSecond']['95%'], data_description_dict['RequestsPerSecond']['99%']]

        analyzed_data = {
            'name' : 'RequestsPerSecond', # used to identify the file name to be saved
            'display_name': 'Requests Per Second', 
            'data_description_chart_data' : {
                'x': x,
                'y': y,
                'text': y, 
                'textposition': 'outside',
                #'marker_color': ['blue', 'yellow', 'red' , 'green', 'orange', 'purple'],
                'name': 'Data Description',
                'xaxis_title_text': 'Percentile',
                'yaxis_title_text': 'Requests Per Second',
            },
        
            
            'average_data_per_second_chart_data': {
                'x': df_copy.index.to_list(),
                'y': df_copy['RequestsPerSecond'].to_list(),
                'mode': 'lines',
                'name': 'Mean Requests Per Second',
                'xaxis_title_text': 'Time (seconds)',
                'yaxis_title_text': 'Requests Per Second',

            }
        }

        fig = self.create_fig(analyzed_data)

        # html summary
        html_div = f"""
                    <div class='card rounded-xl m-10 p-10 border-2'>
                        <p class="text-lg font-bold italic">Summary Information</p>
                        <p class="italic">Mean Requests Per Second: <b>{data_description_dict['RequestsPerSecond']['mean']}</b></p>
                        <p class="italic">Standard Deviation: <b>{data_description_dict['RequestsPerSecond']['std']}</b></p>
                        <p class="italic">25th Percentile: <b>{data_description_dict['RequestsPerSecond']['25%']}</b></p>
                        <p class="italic">50th Percentile: <b>{data_description_dict['RequestsPerSecond']['50%']}</b></p>
                        <p class="italic">75th Percentile: <b>{data_description_dict['RequestsPerSecond']['75%']}</b></p>
                        <p class="italic">90th Percentile: <b>{data_description_dict['RequestsPerSecond']['90%']}</b></p>
                        <p class="italic">95th Percentile: <b>{data_description_dict['RequestsPerSecond']['95%']}</b></p>
                        <p class="italic">99th Percentile: <b>{data_description_dict['RequestsPerSecond']['99%']}</b></p>
                        <p class="italic">Request with lowest latency: <b>{data_description_dict['RequestsPerSecond']['min']}</b></p>
                        <p class="italic">Request with highest latency: <b>{data_description_dict['RequestsPerSecond']['max']}</b></p>


                    </div>
                    """
        self.graphs.append({
            'title': "Requests Per Second Analysis",
            'description': test_description,
            'html_summary': html_div,
            'fig': fig,
            'analyzed_data': analyzed_data
        })
    def create_test_report_html(self):
        # create html file with the report data
        with open(self.data_path + 'report.html', 'w') as f:
            # write the report data to the file
            f.write("<html><head><title>" + self.title + "</title>")
            # use tailwind css
            f.write("<script src='https://cdn.tailwindcss.com'></script>")
            f.write("</head><body>")
            f.write("<div class='card rounded-xl m-10 p-10 border-2'>")
            f.write("<p class='text-2xl font-bold italic'>" + self.title + "</p>")
            f.write("<p class='text-lg font-bold italic'>Test Start Time: <b>" + self.test_start_time + "</b></p>")
            f.write("<p class='text-lg font-bold italic'>Server URL: <b>" + self.server_url + "</b></p>")
            f.write("<p class='text-lg font-bold italic'>Test Description:</p>")
            f.write("<p class='text-lg italic'>" + self.description + "</p>")
            f.write("<p class='font-bold italic'>Run Information:</p>")
            f.write("<p class='italic'>Test Mode: <b>" + self.test_mode + "</b></p>")
            f.write("<p class='italic'>Concurrent Requests: <b>" + str(self.concurrent_requests) + "</b></p>")
            f.write("<p class='italic'>Test Duration: <b>" + str(self.test_duration) + "</b> seconds</p>")
            f.write("<p class='italic'>Total Requests: <b>" + str(self.total_requests) + "</b></p>")
            f.write("<p class='italic'>Successful Requests: <b>" + str(self.successful_requests) + "</b></p>")
            f.write("<p class='italic'>Failed Requests: <b>" + str(self.failed_requests) + "</b></p>")
            f.write("<p class='italic'>Average Requests Per Second (Total Requests / Test Duration): <b>" + str(self.requests_per_second) + "</b></p>")
            f.write("</div>")
            f.write("<hr class='!border-t-4'>")
            for graph in self.graphs:
                f.write("<div class='card rounded-xl m-10 p-10 border-2'>")
                f.write("<p class='text-xl font-bold italic'>" + graph['title'] + "</p>")
                if graph['description'] != '':
                    f.write("<div class='card justify-center rounded-xl m-10 p-10 border-2'>")
                    f.write("<p class='text-lg font-bold italic'> Test Description:</p>")
                    f.write("<p class='text-lg italic'>" + graph['description'] + "</p>")
                    f.write("</div>")
                f.write(graph['html_summary'])
                f.write("<div class='card flex justify-center rounded-xl m-10 p-10 border-2'>")
                f.write(graph['fig'].to_html(full_html=False, include_plotlyjs='cdn'))
                f.write("</div>")
                f.write("</div>") 
                
                f.write("<hr class='!border-t-4'>")


            f.write("</body></html>")
            f.close()
        # save the report data to a json file
        for graph in self.graphs:
            # save analyzed data
            with open(self.data_path + graph['analyzed_data']['name'] + '.json', 'w') as f:
                json.dump(graph['analyzed_data'], f)
                f.close()

            # print("Report created successfully")
    def create_comparison_fig(self, chart_data_list):
        fig = make_subplots(rows=1, cols=2, subplot_titles=("Mean " + chart_data_list[0]["display_name"] +  " Per Second", "Summary Information"))

        # left figure
        for i, chart_data in enumerate(chart_data_list):
            fig.add_trace(go.Scatter(x=chart_data["average_data_per_second_chart_data"]["x"],
                                    y=chart_data["average_data_per_second_chart_data"]["y"],
                                    mode=chart_data["average_data_per_second_chart_data"]["mode"],
                                    name=chart_data["average_data_per_second_chart_data"]["name"],
                                    line=dict(color=['blue', 'red'][i], width=2)),
                        row=1, col=1)
        fig.update_xaxes(title_text=chart_data_list[0]["average_data_per_second_chart_data"]["xaxis_title_text"] , row=1, col=1)
        fig.update_yaxes(title_text=chart_data_list[0]["average_data_per_second_chart_data"]["yaxis_title_text"], row=1, col=1)

        # right figure
        for i, chart_data in enumerate(chart_data_list):
            fig.add_trace(go.Bar(
                            x=chart_data["data_description_chart_data"]["x"],
                            y=chart_data["data_description_chart_data"]["y"],
                            text=chart_data["data_description_chart_data"]["text"],
                            textposition='outside',
                            marker_color=['blue', 'red'][i],
                            name=chart_data["data_description_chart_data"]["name"],
                            legendgroup=chart_data["data_description_chart_data"]["name"],),
                        row=1, col=2)
        fig.update_xaxes(title_text=chart_data_list[0]["data_description_chart_data"]["xaxis_title_text"], row=1, col=2)
        fig.update_yaxes(title_text=chart_data_list[0]["data_description_chart_data"]["yaxis_title_text"], row=1, col=2)


        # add annotation in the footer
        fig.add_annotation(text="Comparison of " + chart_data_list[0]["display_name"] + " Analysis", showarrow=False,
                           xref="paper", yref="paper", x=0.5, y=-0.1)

        # setup the layout for the subplots
        fig.update_layout(title_text=chart_data_list[0]["display_name"] + " Comparison " + self.test_1_name + " vs " + self.test_2_name 
                          , showlegend=True)
        # update the legend title
        fig.update_layout(legend_title_text='Legend')
        

        fig.update_layout(height=600, width=1500)
        # save the figure as png
        #fig.write_image("comparisons/" + chart_data_list[0]["name"] + '.png')
        return fig
    def compare_tests(self, analyses_file_names):
       print("Comparing tests: " + self.test_1_name + " and " + self.test_2_name)
       figs = []
       # loop through the analyses files
       for analysis_file_name in analyses_file_names:
           # load the data for each test
           test_1_data = json.load(open(self.test_1_data_path + analysis_file_name + '.json'))
           test_2_data = json.load(open(self.test_2_data_path + analysis_file_name + '.json'))
           # change the field test_x_data["data_description_chart_data"]["name"] to the test name + the actual name
           test_1_data["data_description_chart_data"]["name"] = self.test_1_name + " " + test_1_data["data_description_chart_data"]["name"]
           test_1_data["average_data_per_second_chart_data"]["name"] = self.test_1_name + " " + test_1_data["average_data_per_second_chart_data"]["name"]
           test_2_data["data_description_chart_data"]["name"] = self.test_2_name + " " + test_2_data["data_description_chart_data"]["name"]
           test_2_data["average_data_per_second_chart_data"]["name"] = self.test_2_name + " " + test_2_data["average_data_per_second_chart_data"]["name"]
           figs.append(self.create_comparison_fig([test_1_data, test_2_data]))
        
            # create html file with the report data
           with open('comparisons/' + self.test_1_name + '_vs_' + self.test_2_name + '_comparison_report.html', 'w') as f:
                    # write the report data to the file
                    f.write("<html><head><title>" + self.test_1_name + " vs " + self.test_2_name + " Comparison Report</title>")
                    # use tailwind css
                    f.write("<script src='https://cdn.tailwindcss.com'></script>")
                    f.write("</head><body>")
                    f.write("<div class='card rounded-xl m-10 p-10 border-2'>")
                    f.write("<p class='text-2xl font-bold italic'>" + self.test_1_name + " vs " + self.test_2_name + " Comparison Report</p>")
                
                    f.write("</div>")
                    f.write("<hr class='!border-t-4'>")
                    for fig in figs:
                        f.write("<div class='card flex justify-center rounded-xl m-10 p-10 border-2'>")
                        f.write(fig.to_html(full_html=False, include_plotlyjs='cdn'))
                        f.write("</div>")
                        f.write("<hr class='!border-t-4'>")
                    f.write("</body></html>")
                    f.close()
                    print("Comparison report created successfully")
       
          
        
       


 




