package com.greencompass.feature.profile

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.KeyboardArrowRight
import androidx.compose.material.icons.outlined.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors

data class ProfileMenuItem(val title: String, val subtitle: String, val icon: ImageVector)

@Composable
fun ProfileScreen(onOpenPlaces: () -> Unit) {
    val menuItems = listOf(
        ProfileMenuItem("Saved places", "4 places", Icons.Outlined.Place),
        ProfileMenuItem("Interests", "3 selected", Icons.Outlined.Star),
        ProfileMenuItem("Update history", "12 updates", Icons.Outlined.History),
        ProfileMenuItem("Settings", "", Icons.Outlined.Settings),
        ProfileMenuItem("Help & support", "", Icons.Outlined.Help),
        ProfileMenuItem("About Green Compass", "", Icons.Outlined.Info)
    )

    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        LazyColumn(modifier = Modifier.fillMaxSize()) {
            item {
                Row(modifier = Modifier.fillMaxWidth().padding(24.dp), verticalAlignment = Alignment.CenterVertically) {
                    Surface(shape = androidx.compose.foundation.shape.CircleShape, color = GreenCompassColors.ForestGreen, modifier = Modifier.size(56.dp)) {
                        Box(contentAlignment = Alignment.Center) { Text(text = "JT", color = Color.White, fontSize = 20.sp, fontWeight = FontWeight.Bold) }
                    }
                    Spacer(Modifier.width(16.dp))
                    Column {
                        Text(text = "John Tabu", fontSize = 20.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal)
                        Text(text = "Lower Valley", fontSize = 14.sp, color = GreenCompassColors.MutedText)
                    }
                }
                Divider(color = GreenCompassColors.Stone)
            }

            items(menuItems.size) { index ->
                val item = menuItems[index]
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 24.dp, vertical = 16.dp)
                        .clickable { if (item.title == "Saved places") onOpenPlaces() },
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(item.icon, contentDescription = null, tint = GreenCompassColors.Charcoal, modifier = Modifier.size(24.dp))
                    Spacer(Modifier.width(16.dp))
                    Column(modifier = Modifier.weight(1f)) {
                        Text(text = item.title, fontSize = 16.sp, color = GreenCompassColors.Charcoal)
                        if (item.subtitle.isNotEmpty()) Text(text = item.subtitle, fontSize = 13.sp, color = GreenCompassColors.MutedText)
                    }
                    Icon(Icons.Default.KeyboardArrowRight, contentDescription = null, tint = GreenCompassColors.MutedText)
                }
                if (index < menuItems.size - 1) Divider(modifier = Modifier.padding(horizontal = 24.dp), color = GreenCompassColors.Stone)
            }
            
            item {
                Spacer(Modifier.height(24.dp))
                TextButton(onClick = {}, modifier = Modifier.fillMaxWidth().padding(horizontal = 24.dp)) {
                    Text(text = "Log out", color = Color(0xFFA93226), fontSize = 16.sp, fontWeight = FontWeight.Medium)
                }
            }
        }
    }
}
