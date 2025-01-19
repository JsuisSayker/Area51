import React from 'react';
import { SvgProps } from 'react-native-svg';
import { SvgXml } from 'react-native-svg';

const homeSvg = `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24">
  <path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z"/>
</svg>
`;

const settingSvg = `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24">
  <path d="M 10.490234 2 C 10.011234 2 9.6017656 2.3385938 9.5097656 2.8085938 L 9.1757812 4.5234375 C 8.3550224 4.8338012 7.5961042 5.2674041 6.9296875 5.8144531 L 5.2851562 5.2480469 C 4.8321563 5.0920469 4.33375 5.2793594 4.09375 5.6933594 L 2.5859375 8.3066406 C 2.3469375 8.7216406 2.4339219 9.2485 2.7949219 9.5625 L 4.1132812 10.708984 C 4.0447181 11.130337 4 11.559284 4 12 C 4 12.440716 4.0447181 12.869663 4.1132812 13.291016 L 2.7949219 14.4375 C 2.4339219 14.7515 2.3469375 15.278359 2.5859375 15.693359 L 4.09375 18.306641 C 4.33275 18.721641 4.8321562 18.908906 5.2851562 18.753906 L 6.9296875 18.1875 C 7.5958842 18.734206 8.3553934 19.166339 9.1757812 19.476562 L 9.5097656 21.191406 C 9.6017656 21.661406 10.011234 22 10.490234 22 L 13.509766 22 C 13.988766 22 14.398234 21.661406 14.490234 21.191406 L 14.824219 19.476562 C 15.644978 19.166199 16.403896 18.732596 17.070312 18.185547 L 18.714844 18.751953 C 19.167844 18.907953 19.66625 18.721641 19.90625 18.306641 L 21.414062 15.691406 C 21.653063 15.276406 21.566078 14.7515 21.205078 14.4375 L 19.886719 13.291016 C 19.955282 12.869663 20 12.440716 20 12 C 20 11.559284 19.955282 11.130337 19.886719 10.708984 L 21.205078 9.5625 C 21.566078 9.2485 21.653063 8.7216406 21.414062 8.3066406 L 19.90625 5.6933594 C 19.66725 5.2783594 19.167844 5.0910937 18.714844 5.2460938 L 17.070312 5.8125 C 16.404116 5.2657937 15.644607 4.8336609 14.824219 4.5234375 L 14.490234 2.8085938 C 14.398234 2.3385937 13.988766 2 13.509766 2 L 10.490234 2 z M 12 8 C 14.209 8 16 9.791 16 12 C 16 14.209 14.209 16 12 16 C 9.791 16 8 14.209 8 12 C 8 9.791 9.791 8 12 8 z"/>
</svg>
`;

const serviceSvg = `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24">
  <path d="M.101,4C.566,1.721,2.586,0,5,0h14c2.414,0,4.434,1.721,4.899,4H.101Zm23.899,2v13c0,2.757-2.243,5-5,5H5c-2.757,0-5-2.243-5-5V6H24ZM3,10c0,.552,.447,1,1,1h5c.553,0,1-.448,1-1s-.447-1-1-1H4c-.553,0-1,.448-1,1Zm12,10c0-.553-.447-1-1-1H4c-.553,0-1,.447-1,1s.447,1,1,1H14c.553,0,1-.447,1-1Zm0-5c0-.553-.447-1-1-1H7c-.553,0-1,.447-1,1s.447,1,1,1h7c.553,0,1-.447,1-1Zm6,5c0-.738-.405-1.376-1-1.723v-6.277c0-1.654-1.346-3-3-3h-1.277c-.346-.595-.984-1-1.723-1-1.105,0-2,.895-2,2s.895,2,2,2c.738,0,1.376-.405,1.723-1h1.277c.552,0,1,.449,1,1v6.277c-.595,.346-1,.984-1,1.723,0,1.105,.895,2,2,2s2-.895,2-2Z"/>
</svg>
`;

const workflowSvg = `
<svg viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
  <g data-name="Layer 2" id="Layer_2">
    <polygon points="8 17.48 12 14.28 16 17.48 16 22 22 22 22 16 17.35 16 13 12.52 13 8 15 8 15 2 9 2 9 8 11 8 11 12.52 6.65 16 2 16 2 22 8 22 8 17.48"/>
  </g>
</svg>
`;

const loginSvg = `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 512 512" width="24" height="24">
  <g>
    <circle cx="256" cy="128" r="128"/>
    <path d="M256,298.667c-105.99,0.118-191.882,86.01-192,192C64,502.449,73.551,512,85.333,512h341.333c11.782,0,21.333-9.551,21.333-21.333C447.882,384.677,361.99,298.784,256,298.667z"/>
  </g>
</svg>
`;

const registerSvg = `
<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24">
  <path d="M15.829,2c-.413-1.164-1.525-2-2.829-2h-2c-1.304,0-2.416,.836-2.829,2H3V21c0,1.654,1.346,3,3,3h12c1.654,0,3-1.346,3-3V2h-5.171Zm-6.829,16h-2v-2h2v2Zm0-4h-2v-2h2v2Zm0-4h-2v-2h2v2Zm8,8h-6v-2h6v2Zm0-4h-6v-2h6v2Zm0-4h-6v-2h6v2Z"/>
</svg>
`;

export const HomeIcon = (props: SvgProps) => (
  <SvgXml xml={homeSvg} {...props} />
);
export const SettingIcon = (props: SvgProps) => (
  <SvgXml xml={settingSvg} {...props} />
);
export const ServiceIcon = (props: SvgProps) => (
  <SvgXml xml={serviceSvg} {...props} />
);
export const WorkflowIcon = (props: SvgProps) => (
  <SvgXml xml={workflowSvg} {...props} />
);
export const LoginIcon = (props: SvgProps) => (
  <SvgXml xml={loginSvg} {...props} />
);
export const RegisterIcon = (props: SvgProps) => (
  <SvgXml xml={registerSvg} {...props} />
);
