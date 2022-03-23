const data = [
  [
    "Name",
    "Date",
    "Price A",
    "Price B",
    "Units A",
    "Units B",
    "Value",
    "HODL",
    "Only A",
    "Only B",
    "APR"
  ],
  [
    "LUNA/UST LP",
    "2022-03-24",
    102.5,
    1,
    48.78048780487805,
    5000,
    10000,
    10000,
    10000,
    10000,
    100
  ],
  [
    "LUNA/OSMO LP",
    "2022-03-24",
    102.5,
    8.333333333333334,
    48.78048780487805,
    600,
    10000,
    10000,
    10000,
    10000,
    125
  ],
  [
    "LUNA/UST LP",
    "2022-03-25",
    105,
    1,
    48.320391170132964,
    5073.641072863961,
    10147.282145727922,
    10121.951219512195,
    10243.90243902439,
    10000,
    94
  ],
  [
    "LUNA/OSMO LP",
    "2022-03-25",
    105,
    8.666666666666668,
    49.310315091140446,
    597.4134328349706,
    10355.166169139493,
    10321.951219512197,
    10243.90243902439,
    10400.000000000002,
    118.5
  ],
  [
    "LUNA/UST LP",
    "2022-03-26",
    107.5,
    1,
    47.87035682901699,
    5146.063359119326,
    10292.126718238653,
    10243.90243902439,
    10487.80487804878,
    10000,
    88
  ],
  [
    "LUNA/OSMO LP",
    "2022-03-26",
    107.5,
    9.000000000000002,
    49.81429546865191,
    595.0040847644533,
    10710.073525760161,
    10643.90243902439,
    10487.80487804878,
    10800.000000000002,
    112
  ],
  [
    "LUNA/UST LP",
    "2022-03-27",
    110,
    1,
    47.429564313160796,
    5217.252074447688,
    10434.504148895376,
    10365.853658536585,
    10731.707317073171,
    10000,
    82
  ],
  [
    "LUNA/OSMO LP",
    "2022-03-27",
    110,
    9.333333333333336,
    50.29357348205806,
    592.7456874671127,
    11064.586166052773,
    10965.853658536587,
    10731.707317073171,
    11200.000000000004,
    105.5
  ],
  [
    "LUNA/UST LP",
    "2022-03-28",
    112.5,
    1,
    46.99726239287874,
    5287.192019198858,
    10574.384038397717,
    10487.80487804878,
    10975.609756097561,
    10000,
    76
  ],
  [
    "LUNA/OSMO LP",
    "2022-03-28",
    112.5,
    9.66666666666667,
    50.74916650234277,
    590.6152998117475,
    11418.562463027123,
    11287.804878048782,
    10975.609756097561,
    11600.000000000004,
    99
  ],
  [
    "LUNA/UST LP",
    "2022-03-29",
    115,
    1,
    46.57276234760598,
    5355.867669974688,
    10711.735339949377,
    10609.756097560976,
    11219.512195121952,
    10000,
    70
  ],
  [
    "LUNA/OSMO LP",
    "2022-03-29",
    115,
    10.000000000000004,
    51.18198207202153,
    588.5927938282474,
    11771.855876564952,
    11609.756097560978,
    11219.512195121952,
    12000.000000000004,
    92.5
  ],
  [
    "LUNA/UST LP",
    "2022-03-30",
    117.5,
    1,
    46.155432032096236,
    5423.263263771308,
    10846.526527542615,
    10731.707317073171,
    11463.414634146342,
    10000,
    64
  ],
  [
    "LUNA/OSMO LP",
    "2022-03-30",
    117.5,
    10.333333333333337,
    51.592832574033565,
    586.6604349144136,
    12124.315654897888,
    11931.707317073175,
    11463.414634146342,
    12400.000000000005,
    86
  ],
  [
    "LUNA/UST LP",
    "2022-03-31",
    120,
    1,
    45.744690629224316,
    5489.362875506917,
    10978.725751013833,
    10853.658536585366,
    11707.317073170732,
    10000,
    58
  ],
  [
    "LUNA/OSMO LP",
    "2022-03-31",
    120,
    10.666666666666671,
    51.982447701744604,
    584.8025366446265,
    12475.787448418705,
    12253.658536585368,
    11707.317073170732,
    12800.000000000005,
    79.5
  ],
  [
    "LUNA/UST LP",
    "2022-04-01",
    122.5,
    1,
    49.42163665016947,
    6054.15048964576,
    12108.30097929152,
    11975.609756097561,
    11951.219512195123,
    10000,
    52
  ],
  [
    "LUNA/OSMO LP",
    "2022-04-01",
    122.5,
    11.000000000000005,
    52.351485111070495,
    583.0051751005576,
    12826.113852212271,
    12575.609756097565,
    11951.219512195123,
    13200.000000000007,
    73
  ],
  [
    "LUNA/UST LP",
    "2022-04-02",
    125,
    1,
    48.98658298629418,
    6123.322873286773,
    12246.645746573546,
    12107.76505724241,
    12195.121951219513,
    10000,
    46
  ],
  [
    "LUNA/OSMO LP",
    "2022-04-02",
    125,
    11.33333333333334,
    52.70053956095547,
    581.2559510399498,
    13175.134890238867,
    12897.56097560976,
    12195.121951219513,
    13600.000000000007,
    66.5
  ],
  [
    "LUNA/UST LP",
    "2022-04-03",
    127.5,
    1,
    48.55709979068318,
    6191.030223312106,
    12382.060446624211,
    12239.920358387259,
    12439.024390243903,
    10000,
    40
  ],
  [
    "LUNA/OSMO LP",
    "2022-04-03",
    127.5,
    11.666666666666673,
    53.03015078996978,
    579.5437907760979,
    13522.688451442293,
    13219.512195121955,
    12439.024390243903,
    14000.000000000007,
    60
  ]
]

export default data
