import React from "react";
import { connect } from "react-redux";
import { Grid, Card, CardContent, Typography, Tooltip } from "@mui/material";
import makeStyles from "@mui/styles/makeStyles";
import moment from "moment";
import { MoneyDirective } from "../../utils/Money";

const useStyles = makeStyles(() => ({
  statCard: {
    background: "#ffffff",
    border: "1px solid #e2e8f0",
    borderRadius: "8px",
    boxShadow: "0 1px 3px rgba(0, 0, 0, 0.06)",
    transition: "all 0.2s ease",
    cursor: "default",
    "&:hover": {
      boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)",
      transform: "translateY(-1px)",
      borderColor: "#cbd5e0",
    },
  },
  statContent: {
    padding: "20px 16px !important",
    textAlign: "center",
  },
  statTitle: {
    fontSize: "0.85rem",
    fontWeight: "800",
    color: "#DC143C",
    marginBottom: "12px",
    textTransform: "uppercase",
    letterSpacing: "0.5px",
  },
  statValue: {
    fontSize: "2.5rem",
    fontWeight: "700",
    fontFamily: "MuseoModerno",
    lineHeight: "1.1",
    minHeight: "60px",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },
  monthlyValue: {
    color: "#059669",
  },
  dailyValue: {
    color: "#7c3aed",
  },
  unusedValue: {
    color: "#dc2626",
  },
  optimizationsValue: {
    color: "#1e40af",
  },
}));

const StatisticsBar = ({ resources }) => {
  const classes = useStyles();

  const collectors = Object.values(resources || {});

  const totalSpent = collectors.reduce((sum, collector) => {
    return sum + (collector.TotalSpent || 0);
  }, 0);

  const dailySpent = totalSpent / 30;

  const totalUnusedResources = collectors
    .filter((collector) => 
      collector.Category === "unused_resource" || 
      (!collector.TotalSpent || collector.TotalSpent === 0)
    )
    .reduce((sum, collector) => sum + (collector.ResourceCount || 0), 0);

  const totalOptimizations = collectors
    .filter((collector) => 
      collector.Category === "potential_cost_saving" || 
      (collector.TotalSpent && collector.TotalSpent > 0)
    )
    .reduce((sum, collector) => sum + (collector.ResourceCount || 0), 0);

  const statistics = [
    {
      title: "💰 Monthly potential savings",
      value: MoneyDirective(totalSpent),
      tooltip: "Total monthly cost that could be saved by optimizing identified resources",
      className: classes.monthlyValue,
    },
    {
      title: "📅 Daily potential savings", 
      value: MoneyDirective(dailySpent),
      tooltip: "Daily cost savings you can achieve by optimizing unused resources",
      className: classes.dailyValue,
    },
    {
      title: "🎯 Cost Optimizations",
      value: totalOptimizations.toLocaleString(),
      tooltip: "Number of resources with potential cost savings",
      className: classes.optimizationsValue,
    },
    {
      title: "🗑️ Unused Resources",
      value: totalUnusedResources.toLocaleString(),
      tooltip: "Number of resources without costs that can be removed",
      className: classes.unusedValue,
    },
  ];

  return (
    <Grid container spacing={2} style={{ marginBottom: "24px" }}>
      {statistics.map((stat, index) => (
        <Grid item xs={12} sm={6} md={3} key={index}>
          <Tooltip title={stat.tooltip} arrow>
            <Card className={classes.statCard}>
              <CardContent className={classes.statContent}>
                <Typography className={classes.statTitle}>
                  {stat.title}
                </Typography>
                <Typography 
                  className={`${classes.statValue} ${stat.className}`}
                  style={{
                    fontSize: "48px",
                    fontWeight: "700",
                    fontFamily: "MuseoModerno",
                    minHeight: "70px",
                  }}
                >
                  {stat.value}
                </Typography>
              </CardContent>
            </Card>
          </Tooltip>
        </Grid>
      ))}
    </Grid>
  );
};

const mapStateToProps = (state) => ({
  resources: state.resources.resources,
});

export default connect(mapStateToProps)(StatisticsBar);
