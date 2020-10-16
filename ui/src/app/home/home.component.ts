import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'tide-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss'],
})
export class HomeComponent implements OnInit {

  constructor() { }

  barData = [
    { data: this.genRandomData(5, 100), label: 'CPU%' },
    { data: this.genRandomData(5, 100), label: 'Memory%' },
    { data: this.genRandomData(5, 100), label: 'Disk%' },
  ]

  barLabels = [
    'Folding@home',
    'CAS@home',
    'Asteroids@home',
    'BOINC@TACC',
    'LHC@home@home',
  ]

  barOptions = {
    title: {
      display: true,
      text: 'Percent of Resource Usage for Each Project'
    }
  }

  lineLabels = [
    'January',
    'February',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December'
  ];

  lineOptions = {
    title: {
      display: true,
      text: 'Running Hours for Each Project',
    }
  };

  lineData = [
    { data: this.genRandomData(12, 220), label: 'Folding@home' },
    { data: this.genRandomData(12, 220), label: 'CAS@home' },
    { data: this.genRandomData(12, 720), label: 'Asteroids@home' },
    { data: this.genRandomData(12, 720), label: 'BOINC@TACC' },
    { data: this.genRandomData(12, 720), label: 'LHC@home@home' },
  ];

  bubbleData = [
    {
      data: new Array(500)
      .fill({})
      .map(() => ({
        x: Math.random() * 60,
        y: Math.random() * 100,
        r: Math.random() * 10,
      })),
      label: 'VM',
    }
  ];

  bubbleOptions = {
    title: {
      display: true,
      text: 'Resource Distribution for Each VM',
    }
  };

  ngOnInit(): void {

  }

  genRandomData(size, range) {
    return new Array(size).fill('').map(a => Math.random() * range)
  }

}
