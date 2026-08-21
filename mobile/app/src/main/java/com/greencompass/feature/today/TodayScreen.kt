package com.greencompass.feature.today

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.outlined.WaterDrop
import androidx.compose.material.icons.outlined.Cloud
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.greencompass.core.ui.GreenCompassColors
import com.greencompass.core.ui.LoadingSkeleton

@Composable
fun TodayRoute(viewModel: TodayViewModel = hiltViewModel()) {
    val state by viewModel.uiState.collectAsState()
    TodayScreen(state = state, onRefresh = viewModel::refresh)
}

@Composable
fun TodayScreen(state: TodayUiState, onRefresh: () -> Unit) {
    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        LazyColumn(
            modifier = Modifier.fillMaxSize().padding(24.dp),
            verticalArrangement = Arrangement.spacedBy(20.dp)
        ) {
            item {
                if (state.isLoading) {
                    LoadingSkeleton(modifier = Modifier.fillMaxWidth().height(40.dp))
                    Spacer(Modifier.height(16.dp))
                    LoadingSkeleton(modifier = Modifier.fillMaxWidth(0.6f).height(32.dp))
                } else {
                    Text(text = "Good morning", fontSize = 14.sp, color = GreenCompassColors.MutedText)
                    Text(text = "Lower Valley", fontSize = 28.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal)
                }
            }

            // Weather/Conditions Card
            item {
                if (state.isLoading) {
                    LoadingSkeleton(modifier = Modifier.fillMaxWidth().height(120.dp))
                } else {
                    Card(
                        shape = RoundedCornerShape(16.dp),
                        colors = CardDefaults.cardColors(containerColor = GreenCompassColors.SoftSage)
                    ) {
                        Row(modifier = Modifier.padding(20.dp), verticalAlignment = Alignment.CenterVertically) {
                            Icon(Icons.Outlined.Cloud, contentDescription = null, tint = GreenCompassColors.ForestGreen, modifier = Modifier.size(48.dp))
                            Spacer(Modifier.width(16.dp))
                            Column {
                                Text(text = "24°C", fontSize = 32.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal)
                                Text(text = "Partly cloudy", fontSize = 15.sp, color = GreenCompassColors.MutedText)
                            }
                        }
                    }
                }
            }

            // Key Updates
            item {
                Text(text = "Key updates", fontSize = 16.sp, fontWeight = FontWeight.SemiBold, color = GreenCompassColors.Charcoal)
            }

            if (state.isLoading) {
                item { LoadingSkeleton(modifier = Modifier.fillMaxWidth().height(80.dp)) }
            } else if (state.data?.updates?.isEmpty() == true) {
                item {
                    Card(shape = RoundedCornerShape(12.dp), colors = CardDefaults.cardColors(containerColor = Color.White), border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)) {
                        Row(modifier = Modifier.padding(16.dp), verticalAlignment = Alignment.CenterVertically) {
                            Icon(Icons.Outlined.Cloud, contentDescription = null, tint = GreenCompassColors.ForestGreen)
                            Spacer(Modifier.width(12.dp))
                            Column {
                                Text(text = "No important alerts", fontSize = 15.sp, fontWeight = FontWeight.Medium, color = GreenCompassColors.Charcoal)
                                Text(text = "Everything looks normal today.", fontSize = 13.sp, color = GreenCompassColors.MutedText)
                            }
                        }
                    }
                }
            }

            // Water Conditions
            item {
                if (!state.isLoading) {
                    Card(shape = RoundedCornerShape(12.dp), colors = CardDefaults.cardColors(containerColor = Color.White), border = androidx.compose.foundation.BorderStroke(1.dp, GreenCompassColors.Stone)) {
                        Row(modifier = Modifier.padding(16.dp), verticalAlignment = Alignment.CenterVertically) {
                            Icon(Icons.Outlined.WaterDrop, contentDescription = null, tint = GreenCompassColors.SkyBlue)
                            Spacer(Modifier.width(12.dp))
                            Column {
                                Text(text = "Water conditions", fontSize = 15.sp, fontWeight = FontWeight.Medium, color = GreenCompassColors.Charcoal)
                                Text(text = "Normal · No changes in the last 24h.", fontSize = 13.sp, color = GreenCompassColors.MutedText)
                            }
                        }
                    }
                }
            }
        }
    }
}
