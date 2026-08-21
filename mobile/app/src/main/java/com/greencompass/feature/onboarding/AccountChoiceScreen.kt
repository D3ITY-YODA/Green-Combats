package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors

@Composable
fun AccountChoiceScreen(onPersonal: () -> Unit, onOrganization: () -> Unit, onInvitation: () -> Unit) {
    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(32.dp), horizontalAlignment = Alignment.CenterHorizontally, verticalArrangement = Arrangement.Center) {
            Text(text = "How would you like to continue?", fontSize = 24.sp, fontWeight = FontWeight.Bold, color = GreenCompassColors.Charcoal, textAlign = TextAlign.Center, modifier = Modifier.padding(bottom = 40.dp))
            
            Button(onClick = onPersonal, modifier = Modifier.fillMaxWidth().height(56.dp), colors = ButtonDefaults.buttonColors(containerColor = GreenCompassColors.ForestGreen), shape = RoundedCornerShape(12.dp)) { Text("Use Green Compass personally", fontSize = 16.sp) }
            Spacer(modifier = Modifier.height(16.dp))
            OutlinedButton(onClick = onOrganization, modifier = Modifier.fillMaxWidth().height(56.dp), shape = RoundedCornerShape(12.dp)) { Text("Join an organization", fontSize = 16.sp, color = GreenCompassColors.ForestGreen) }
            Spacer(modifier = Modifier.height(16.dp))
            TextButton(onClick = onInvitation, modifier = Modifier.fillMaxWidth()) { Text("I was invited by an organization", fontSize = 15.sp, color = GreenCompassColors.ForestGreen) }
        }
    }
}
