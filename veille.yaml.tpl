---
roll_windows: [30, 90, 365]
ts_interval: 15
false_positives:
  # There were some false positives during this time
  - start: '2015-08-03 22:00'
    end: '2015-08-04 03:40'
    service_patterns: ['*']
  # This service kept throwing false positives during these weeks
  - start: '2015-04-12 00:00'
    end: '2015-04-26 00:00'
    sevice_patterns:
      - 'somehost;My Glitchy Servicecheck'
